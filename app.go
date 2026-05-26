package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/wailsapp/go-webview2/webviewloader"
	wails_runtime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (m *App) SetContext(ctx context.Context) {
	m.ctx = ctx
}

func (a *App) Cwd() (string, error) {
	if exePath, err := os.Executable(); err != nil {
		return "", err
	} else {
		return filepath.Dir(exePath), nil
	}
}

func (a *App) SelectFolder(relative bool) (string, error) {
	if path, err := wails_runtime.OpenDirectoryDialog(a.ctx, wails_runtime.OpenDialogOptions{}); err != nil || path == "" {
		return "", err
	} else if relative {
		if exePath, err := os.Executable(); err != nil {
			return "", err
		} else {
			return filepath.Rel(filepath.Dir(exePath), path)
		}
	} else {
		return path, nil
	}
}

func (a *App) SelectFile(relative bool) (string, error) {
	if path, err := wails_runtime.OpenFileDialog(a.ctx, wails_runtime.OpenDialogOptions{}); err != nil || path == "" {
		return "", err
	} else if relative {
		if exePath, err := os.Executable(); err != nil {
			return "", err
		} else {
			return filepath.Rel(filepath.Dir(exePath), path)
		}
	} else {
		return path, nil
	}
}

func (a App) PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (a App) ExecutableExists(path string) bool {
	_, err := exec.LookPath(path)
	return err == nil
}

func (a App) WebView2Version() (string, error) {
	return webviewloader.GetAvailableCoreWebView2BrowserVersionString(pathWV2)
}

func (a App) WebView2Path() string {
	return pathWV2
}

func (a App) AppConfigPath() string {
	return dirConf
}

func (a App) AppDriverPath() string {
	return dirDir
}

func (a App) AppVersion() string {
	return version.String()
}

func (a App) AppBinaryType() string {
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x64"
	} else if arch == "386" {
		arch = "x86"
	}
	return fmt.Sprintf("%s-%s", runtime.GOOS, arch)
}

// TriggerNativeUpdate downloads a ZIP from downloadUrl, extracts install-it.exe,
// performs the Windows file-rename trick to replace the running binary, spawns the
// new process, and exits the current instance.
func (a App) TriggerNativeUpdate(downloadUrl string) error {
	// Download the ZIP payload
	resp, err := http.Get(downloadUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("app: update download failed: %s", resp.Status)
	}

	// Write to a temp file inside dirRoot so rename stays on the same volume
	tmpZip, err := os.CreateTemp(dirRoot, "update-*.zip")
	if err != nil {
		return err
	}
	tmpZipPath := tmpZip.Name()
	defer os.Remove(tmpZipPath)

	if _, err := io.Copy(tmpZip, resp.Body); err != nil {
		tmpZip.Close()
		return err
	}
	tmpZip.Close()

	// Extract install-it.exe from the ZIP
	newBinPath := filepath.Join(dirRoot, "install-it.exe.new")
	if err := extractBinaryFromZip(tmpZipPath, newBinPath); err != nil {
		os.Remove(newBinPath)
		return err
	}

	exePath := filepath.Join(dirRoot, "install-it.exe")
	oldBinPath := filepath.Join(dirRoot, "install-it.exe.old")

	// Windows rename trick: rename running exe out of the way
	if err := os.Rename(exePath, oldBinPath); err != nil {
		os.Remove(newBinPath)
		return fmt.Errorf("app: failed to rename current binary: %w", err)
	}

	// Rename new binary into place
	if err := os.Rename(newBinPath, exePath); err != nil {
		// Best-effort restore
		os.Rename(oldBinPath, exePath)
		return fmt.Errorf("app: failed to rename new binary into place: %w", err)
	}

	// Spawn the new process detached from the current one
	process, err := os.StartProcess(exePath, []string{exePath}, &os.ProcAttr{
		Dir:   dirRoot,
		Env:   os.Environ(),
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	})
	if err != nil {
		return fmt.Errorf("app: failed to launch updated binary: %w", err)
	}
	process.Release()

	os.Exit(0)
	return nil // unreachable
}

// extractBinaryFromZip finds install-it.exe inside the ZIP (by base-name suffix)
// and writes it to destPath.
func extractBinaryFromZip(zipPath, destPath string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("app: failed to open update ZIP: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.EqualFold(filepath.Base(f.Name), "install-it.exe") {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			dest, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer dest.Close()

			if _, err := io.Copy(dest, rc); err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("app: install-it.exe not found in update package")
}
