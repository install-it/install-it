package porter

import (
	"archive/zip"
	"errors"
	"fmt"
	"install-it/pkg/status"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ImportPreview describes what's in a ZIP, returned by ValidateZip and DownloadAndValidate.
type ImportPreview struct {
	ExportedAt  time.Time `json:"exportedAt"`
	HasSettings bool      `json:"hasSettings"`
	HasData     bool      `json:"hasData"`
	HasDatabase bool      `json:"hasDatabase"`
	HasDrivers  bool      `json:"hasDrivers"`
	DriverCount int       `json:"driverCount"`
	DriverSize  int64     `json:"driverSize"`
}

// ImportOptions is the user's category selection for import.
type ImportOptions struct {
	Settings bool `json:"settings"`
	Data     bool `json:"data"`
}

// Porter is NOT safe for concurrent use. All methods must be called sequentially
// from a single goroutine (Wails serializes frontend-to-Go calls).
//
// Porter manages the porting process including export, import, and progress tracking.
type Porter struct {
	DirRoot string   // Root directory for import/export operations
	Targets []string // Target directories to be backed up or compressed

	OnBeforeBackup func() error // Called before backup to close DB
	OnAfterImport  func() error // Called after import to reopen DB

	job      *job   // Current job, nil when idle
	tempPath string // Temp file from DownloadAndValidate, cleaned up after ImportFromURL.
	// NOTE: tempPath is not synchronized — assumes single-caller sequential access
	// (DownloadAndValidate → ImportFromURL). No concurrent calls from frontend.
}

func (p *Porter) Status() status.Status {
	if p.job == nil {
		return status.Pending
	}
	p.job.mu.Lock()
	defer p.job.mu.Unlock()
	return p.job.status
}

func (p *Porter) Abort() error {
	if p.job == nil {
		return errors.New("porter: no running porting job")
	}
	p.job.mu.Lock()
	defer p.job.mu.Unlock()
	if p.job.status != status.Running {
		return errors.New("porter: no running porting job")
	}
	p.job.cancel()
	return nil
}

func (p *Porter) Progress() (JobSnapshot, error) {
	if p.job == nil {
		return JobSnapshot{}, errors.New("porter: no started job")
	}
	return p.job.snapshot(), nil
}

func (p *Porter) Export(dest string) (err error) {
	p.job = newJob()
	p.job.start()
	defer func() {
		if err != nil {
			p.job.fail(err)
		} else {
			p.job.complete()
		}
	}()

	return toZip(p.job, dest, p.DirRoot, p.Targets)
}

// ValidateZip reads the manifest and scans ZIP entries for recognized data.
func (p *Porter) ValidateZip(path string) (ImportPreview, error) {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return ImportPreview{}, fmt.Errorf("porter: cannot open zip: %w", err)
	}
	defer zr.Close()

	m, err := readManifest(zr)
	if err != nil {
		return ImportPreview{}, err
	}

	if m.FormatVersion != CurrentFormatVersion {
		return ImportPreview{}, fmt.Errorf("porter: unsupported archive format version %d (expected %d)", m.FormatVersion, CurrentFormatVersion)
	}

	preview := ImportPreview{
		ExportedAt: m.ExportedAt,
	}

	for _, f := range zr.File {
		name := filepath.ToSlash(f.Name)

		switch {
		case name == "conf/setting.json":
			preview.HasSettings = true
		case name == "conf/data.db":
			preview.HasDatabase = true
		case strings.HasPrefix(name, "drivers/"):
			if !preview.HasDrivers {
				preview.HasDrivers = true
			}
			if !f.FileInfo().IsDir() {
				preview.DriverCount++
				preview.DriverSize += int64(f.FileInfo().Size())
			}
		}
	}

	preview.HasData = preview.HasDatabase || preview.HasDrivers

	if !preview.HasSettings && !preview.HasData {
		return ImportPreview{}, fmt.Errorf("porter: no install-it data found in archive")
	}

	return preview, nil
}

// DownloadAndValidate downloads a ZIP from a URL, validates it, and stores the temp path for ImportFromURL.
func (p *Porter) DownloadAndValidate(url string) (preview ImportPreview, err error) {
	if p.tempPath != "" {
		os.Remove(p.tempPath)
		p.tempPath = ""
	}

	p.job = newJob()
	p.job.start()
	defer func() {
		if err != nil {
			if p.tempPath != "" {
				os.Remove(p.tempPath)
				p.tempPath = ""
			}
			p.job.fail(err)
		} else {
			p.job.complete()
		}
	}()

	path, err := download(p.job, url)
	if err != nil {
		return ImportPreview{}, err
	}

	p.tempPath = path

	preview, err = p.ValidateZip(path)
	if err != nil {
		os.Remove(path)
		p.tempPath = ""
		return ImportPreview{}, err
	}

	return preview, nil
}

// ImportFromFile backs up existing files, extracts selected categories, and cleans up or rolls back on failure.
func (p *Porter) ImportFromFile(path string, opts ImportOptions) (err error) {
	p.job = newJob()
	p.job.start()

	var (
		dbClosed  bool
		timestamp string
	)

	defer func() {
		if dbClosed && p.OnAfterImport != nil {
			if reopenErr := p.OnAfterImport(); reopenErr != nil {
				if err == nil {
					err = fmt.Errorf("porter: import succeeded but failed to reopen database: %w", reopenErr)
				} else {
					err = fmt.Errorf("porter: %w (additionally, failed to reopen database: %v)", err, reopenErr)
				}
			}
		}
	}()

	preview, err := p.ValidateZip(path)
	if err != nil {
		p.job.fail(err)
		return err
	}

	// Compute backup set: files and dirs that are in ZIP ∩ selected by opts ∩ exist on disk
	var backupFiles, backupDirs []string

	if opts.Settings {
		if preview.HasSettings {
			if _, statErr := os.Stat(filepath.Join(p.DirRoot, "conf", "setting.json")); statErr == nil {
				backupFiles = append(backupFiles, filepath.Join("conf", "setting.json"))
			}
		}
	}

	if opts.Data {
		if !preview.HasData {
			err := fmt.Errorf("porter: selected categories not found in archive")
			p.job.fail(err)
			return err
		}
		if preview.HasDatabase {
			if _, statErr := os.Stat(filepath.Join(p.DirRoot, "conf", "data.db")); statErr == nil {
				backupFiles = append(backupFiles, filepath.Join("conf", "data.db"))
			}
		}
		if preview.HasDrivers {
			if _, statErr := os.Stat(filepath.Join(p.DirRoot, "drivers")); statErr == nil {
				backupDirs = append(backupDirs, "drivers")
			}
		}
	} else {
		if !opts.Settings || !preview.HasSettings {
			err := fmt.Errorf("porter: nothing to import — no categories selected")
			p.job.fail(err)
			return err
		}
	}

	if len(backupFiles) == 0 && len(backupDirs) == 0 {
		err := fmt.Errorf("porter: nothing to backup or import — selected items do not exist on disk or in archive")
		p.job.fail(err)
		return err
	}

	for _, f := range backupFiles {
		if strings.HasSuffix(f, "conf/data.db") || strings.HasSuffix(f, "conf\\data.db") {
			dbClosed = true
			break
		}
	}

	if dbClosed && p.OnBeforeBackup != nil {
		p.job.msg("Closing database for backup...")
		if err := p.OnBeforeBackup(); err != nil {
			p.job.fail(err)
			return fmt.Errorf("porter: error closing database: %w", err)
		}
	}

	p.job.msg("Backing up existing files...")
	timestamp, err = backup(p.job, p.DirRoot, backupFiles, backupDirs)
	if err != nil {
		p.job.fail(err)
		return err
	}

	p.job.msg("Extracting archive...")
	err = fromZip(p.job, path, p.DirRoot, opts)
	if err != nil {
		p.job.msg("Extraction failed, rolling back...")
		rollbackErr := rollback(p.job, p.DirRoot, timestamp, backupFiles, backupDirs)
		p.job.fail(err)
		if rollbackErr != nil {
			return fmt.Errorf("porter: %w (rollback: %v)", err, rollbackErr)
		}
		return err
	}

	p.job.msg("Cleaning up backups...")
	if err := cleanupBackups(p.job, p.DirRoot, timestamp); err != nil {
		p.job.msg(fmt.Sprintf("Warning: cleanup issue: %v", err))
	}

	p.job.complete()
	return nil
}

// ImportFromURL imports from the temp file left by DownloadAndValidate.
func (p *Porter) ImportFromURL(opts ImportOptions) error {
	if p.tempPath == "" {
		return fmt.Errorf("porter: no downloaded file — call DownloadAndValidate first")
	}

	defer func() {
		os.Remove(p.tempPath)
		p.tempPath = ""
	}()

	return p.ImportFromFile(p.tempPath, opts)
}

// RecoverOrphanedBackups restores leftover .porter-* folders from interrupted imports. Called on startup.
func (p *Porter) RecoverOrphanedBackups() error {
	matches, _ := filepath.Glob(filepath.Join(p.DirRoot, ".porter-*"))
	for _, backupDir := range matches {
		walkErr := filepath.Walk(backupDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || path == backupDir || info.IsDir() {
				return nil
			}
			rel, err := filepath.Rel(backupDir, path)
			if err != nil {
				return err
			}
			original := filepath.Join(p.DirRoot, rel)
			if _, statErr := os.Stat(original); statErr == nil {
				os.Remove(original)
			}
			os.MkdirAll(filepath.Dir(original), os.ModePerm)
			return os.Rename(path, original)
		})

		if walkErr != nil {
			// Leave backup dir for manual recovery or next startup retry
			continue
		}
		os.RemoveAll(backupDir)
	}
	return nil
}
