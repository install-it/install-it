package porter

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Calculates the total size of a directory and its subdirectories.
// If exclDir is true, directory sizes are excluded from the total.
func dirSize(target string, exclDir bool) (int64, error) {
	var size int64
	err := filepath.Walk(target, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() || (!exclDir && info.IsDir()) {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

// Finalises the progress tracker based on the presence of an error.
func updateProgress(progress *Progress, err error) {
	if err != nil {
		if err != context.Canceled {
			progress.message <- err.Error()
		}
		progress.Fail(err)
	} else {
		progress.Complete()
	}
}

// Compresses the specified directories into a single ZIP file at the destination path.
func toZip(tracker *Progress, dest string, directories ...string) (err error) {
	tracker.Start(0)
	defer func() { updateProgress(tracker, err) }()

	for _, dir := range directories {
		if size, err := dirSize(dir, false); err == nil {
			tracker.Total += size
		}
	}

	file, err := os.Create(path.Join(dest, "install-it.zip"))
	if err != nil {
		return err
	}
	defer file.Close()

	zwriter := zip.NewWriter(file)
	defer zwriter.Close()

	for _, dir := range directories {
		err = filepath.Walk(dir, func(filePath string, info os.FileInfo, err error) error {
			if tracker.context.Err() == context.Canceled {
				return tracker.context.Err()
			}
			if err != nil {
				return err
			}

			tracker.message <- fmt.Sprintf("Packing: %s", filePath)

			if info.IsDir() {
				tracker.Accumulate(info.Size())
				return nil
			}

			srcFile, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer srcFile.Close()

			zipEntry, err := zwriter.Create(filePath)
			if err != nil {
				return err
			}

			if _, err = io.Copy(zipEntry, srcFile); err != nil {
				return err
			}

			tracker.Accumulate(info.Size())
			return nil
		})

		if err != nil {
			return err
		}
	}

	tracker.message <- fmt.Sprintf("All files were packed into: %s", file.Name())
	return nil
}

// fromZip extracts a ZIP archive to the specified destination directory.
func fromZip(tracker *Progress, orig string, dest string) (err error) {
	tracker.Start(0)
	defer func() { updateProgress(tracker, err) }()

	zreader, err := zip.OpenReader(orig)
	if err != nil {
		return err
	}
	defer zreader.Close()

	os.MkdirAll(dest, os.ModePerm)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(zf *zip.File) error {
		if tracker.context.Err() == context.Canceled {
			return tracker.context.Err()
		}

		zfreader, err := zf.Open()
		if err != nil {
			return err
		}
		defer zfreader.Close()

		extractPath := filepath.Join(dest, zf.Name)

		// Prevent ZipSlip vulnerability
		if !strings.HasPrefix(extractPath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("porting: illegal file path: %s", extractPath)
		}

		tracker.message <- fmt.Sprintf("Unpacking: %s", zf.Name)

		if zf.FileInfo().IsDir() {
			return os.MkdirAll(extractPath, zf.Mode())
		}

		os.MkdirAll(filepath.Dir(extractPath), zf.Mode())
		outFile, err := os.OpenFile(extractPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zf.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, zfreader)
		return err
	}

	for _, zf := range zreader.File {
		tracker.Total += zf.FileInfo().Size()
	}

	for _, zf := range zreader.File {
		if err := extractAndWriteFile(zf); err != nil {
			return err
		}
		tracker.Accumulate(zf.FileInfo().Size())
	}

	tracker.Complete()
	return nil
}

// Fetches a ZIP file from the given URL and saves it to a temporary file.
func download(tracker *Progress, url string) (path string, err error) {
	tracker.Start(0)
	defer func() { updateProgress(tracker, err) }()

	req, err := http.NewRequestWithContext(tracker.context, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpFile, err := os.CreateTemp("", "*.zip")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	tracker.Total = resp.ContentLength
	tracker.message <- "Downloading..."

	if _, err = io.Copy(tmpFile, io.TeeReader(resp.Body, tracker)); err != nil {
		return "", err
	}

	return filepath.Abs(tmpFile.Name())
}

// Renames each target directory by appending "_old" to create a backup.
func backup(tracker *Progress, targets []string) (err error) {
	tracker.Start(2)
	defer func() { updateProgress(tracker, err) }()

	tracker.message <- "Creating backups..."

	for _, d := range targets {
		if err := os.Rename(d, fmt.Sprintf("%s_old", d)); err != nil {
			return err
		}
		tracker.message <- fmt.Sprintf("%[1]s -> %[1]s_old", d)
		tracker.Accumulate(1)
	}
	return nil
}

// Removes or restores backup directories based on the restore flag.
func cleanup(tracker *Progress, targets []string, restore bool) (err error) {
	tracker.Start(int64(len(targets)))
	defer func() { updateProgress(tracker, err) }()

	if restore {
		tracker.message <- "Restoring backups..."
		for _, d := range targets {
			if err := os.RemoveAll(d); err != nil {
				return err
			}
			if err := os.Rename(fmt.Sprintf("%s_old", d), d); err != nil {
				return err
			}

			tracker.message <- fmt.Sprintf("%[1]s_old -> %[1]s", d)
			tracker.Accumulate(1)
		}
	} else {
		tracker.message <- "Removing backups..."
		for _, d := range targets {
			path := fmt.Sprintf("%s_old", d)
			tracker.message <- fmt.Sprintf("Removing: %s", path)

			if err := os.RemoveAll(path); err != nil {
				tracker.message <- err.Error()
				tracker.message <- fmt.Sprintf("⚠️ Unable to remove backup \"%s\", please consider removing it manually", d)
			} else {
				tracker.Accumulate(1)
			}
		}
	}
	return nil
}
