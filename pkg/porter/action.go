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
//
// exclDir is a boolean flag indicating whether to exclude directories from the size calculation.
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

// Return True if all elements tested by pred is true.
func all[T any](ts []T, pred func(T) bool) bool {
	for _, t := range ts {
		if !pred(t) {
			return false
		}
	}
	return true
}

// Return True if some elements tested by pred is true.
func some[T any](ts []T, pred func(T) bool) bool {
	for _, t := range ts {
		if pred(t) {
			return true
		}
	}
	return false
}

func updateProgress(progress *Progress, err error) {
	if err != nil {
		progress.message <- err.Error()
		progress.Fail(err)
	} else {
		progress.Complete()
	}
}

func toZip(tracker *Progress, dest string, directories ...string) (err error) {
	var total int64 = 0
	for _, directory := range directories {
		if size, err := dirSize(directory, false); err == nil {
			total += size
		}
	}

	tracker.Start(total)
	defer updateProgress(tracker, err)

	file, err := os.Create(path.Join(dest, "driver-box.zip"))
	if err != nil {
		tracker.message <- err.Error()
		tracker.Fail(err)
		return err
	}
	defer file.Close()

	zwriter := zip.NewWriter(file)
	defer zwriter.Close()

	for _, path := range directories {
		err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if tracker.context.Err() == context.Canceled {
				return tracker.context.Err()
			}

			if err != nil {
				return err
			}

			tracker.message <- fmt.Sprintf("Packing: %s", path)

			if info.IsDir() {
				tracker.Accumulate(info.Size())
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			f, err := zwriter.Create(path)
			if err != nil {
				return err
			}

			_, err = io.Copy(f, file)
			if err != nil {
				return err
			}

			tracker.Accumulate(info.Size())
			return nil
		})

		if err != nil {
			tracker.message <- err.Error()
			tracker.Fail(err)
			return err
		}
	}

	tracker.message <- fmt.Sprintf("All files were packed into: %s", file.Name())
	return nil
}

// Reference: https://stackoverflow.com/a/24792688
func fromZip(tracker *Progress, orig string, dest string) error {
	zreader, err := zip.OpenReader(orig)
	if err != nil {
		tracker.Fail(err)
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

		path := filepath.Join(dest, zf.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("porting: illegal file path: %s", path)
		}

		tracker.message <- fmt.Sprintf("Unpacking: %s", zf.Name)

		if zf.FileInfo().IsDir() {
			os.MkdirAll(path, zf.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), zf.Mode())

			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zf.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, zfreader)
			if err != nil {
				return err
			}
		}
		return nil
	}

	var total int64 = 0
	for _, zf := range zreader.File {
		total += zf.FileInfo().Size()
	}

	tracker.Start(total)

	for _, f := range zreader.File {
		if err := extractAndWriteFile(f); err != nil {
			tracker.Fail(err)
			return err
		}
		tracker.Accumulate(f.FileInfo().Size())
	}

	tracker.Complete()
	return nil
}

func download(tracker *Progress, url string) (path string, err error) {
	tracker.Start(1) // placeholder value before establish connection to URL
	defer updateProgress(tracker, err)

	request, err := http.NewRequestWithContext(tracker.context, "GET", url, nil)
	if err != nil {
		tracker.Fail(err)
		return "", err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		tracker.Fail(err)
		return "", err
	}
	defer response.Body.Close()

	file, err := os.CreateTemp("", "*.zip")
	if err != nil {
		tracker.Fail(err)
		return "", err
	}
	defer file.Close()

	tracker.Start(response.ContentLength)
	tracker.message <- "Downloading..."

	if _, err = io.Copy(file, io.TeeReader(response.Body, tracker)); err != nil {
		tracker.Fail(err)
		return "", err
	}

	return filepath.Abs(file.Name())
}

func backup(tracker *Progress, targets []string) (err error) {
	tracker.Start(2)
	defer updateProgress(tracker, err)

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

func cleanup(tracker *Progress, targets []string, restore bool) (err error) {
	tracker.Start(2)
	defer updateProgress(tracker, err)

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
		return nil
	} else {
		tracker.message <- "Cleaning up backups..."

		for _, d := range targets {
			tracker.message <- fmt.Sprintf("Removing: %s_old", d)
			if err := os.RemoveAll(fmt.Sprintf("%s_old", d)); err != nil {
				// not able to removing backup is not a critical problem
				tracker.message <- err.Error()
			} else {
				tracker.Accumulate(1)
			}
		}
		return nil
	}
}
