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

/*
Return True if all elements tested by pred is true.
*/
func all[T any](ts []T, pred func(T) bool) bool {
	for _, t := range ts {
		if !pred(t) {
			return false
		}
	}
	return true
}

/*
Return True if some elements tested by pred is true.
*/
func some[T any](ts []T, pred func(T) bool) bool {
	for _, t := range ts {
		if pred(t) {
			return true
		}
	}
	return false
}

func toZip(tracker *ProgressTracker, dest string, directories ...string) error {
	var total int64 = 0
	for _, directory := range directories {
		if size, err := dirSize(directory, false); err == nil {
			total += size
		}
	}

	tracker.Start("compression", total)

	file, err := os.Create(path.Join(dest, "driver-box.zip"))
	if err != nil {
		tracker.messages <- err.Error()
		tracker.Fail("compression", err)
		return err
	}
	defer file.Close()

	zwriter := zip.NewWriter(file)
	defer zwriter.Close()

	for _, path := range directories {
		err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if tracker.ctx.Err() == context.Canceled {
				return tracker.ctx.Err()
			}

			if err != nil {
				return err
			}

			tracker.messages <- fmt.Sprintf("Packing: %s", path)

			if info.IsDir() {
				tracker.Accumulate("compression", info.Size())
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

			tracker.Accumulate("compression", info.Size())
			return nil
		})

		if err != nil {
			tracker.messages <- err.Error()
			tracker.Fail("compression", err)
			return err
		}
	}

	tracker.messages <- fmt.Sprintf("All files were packed into: %s", file.Name())

	tracker.Complete("compression")
	return nil
}

// Reference: https://stackoverflow.com/a/24792688
func fromZip(tracker *ProgressTracker, orig string, dest string) error {
	zreader, err := zip.OpenReader(orig)
	if err != nil {
		tracker.Fail("decompression", err)
		return err
	}

	defer zreader.Close()

	os.MkdirAll(dest, os.ModePerm)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(zf *zip.File) error {
		if tracker.ctx.Err() == context.Canceled {
			return tracker.ctx.Err()
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

		tracker.messages <- fmt.Sprintf("Unpacking: %s", zf.Name)

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

	tracker.Start("decompression", total)

	for _, f := range zreader.File {
		if err := extractAndWriteFile(f); err != nil {
			tracker.Fail("decompression", err)
			return err
		}
		tracker.Accumulate("decompression", f.FileInfo().Size())
	}

	tracker.Complete("decompression")
	return nil
}

func download(tracker *ProgressTracker, url string, dest string) error {
	tracker.Start("download", 1) // placeholder value before establish connection to URL

	request, err := http.NewRequestWithContext(tracker.ctx, "GET", url, nil)
	if err != nil {
		tracker.Fail("download", err)
		return err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		tracker.Fail("download", err)
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(dest)
	if err != nil {
		tracker.Fail("download", err)
		return err
	}
	defer file.Close()

	tracker.Start("download", response.ContentLength)
	tracker.messages <- "Downloading..."

	if progress, ok := tracker.progs["download"]; ok {
		if _, err = io.Copy(file, io.TeeReader(response.Body, progress)); err != nil {
			tracker.Fail("download", err)
			return err
		}
	} else {
		// panic(err)
	}

	tracker.Complete("download")
	return nil
}

func (p Porter) backup() error {
	p.tracker.Start("backup", 2)

	p.tracker.messages <- "Creating backups..."

	for _, d := range []string{p.DirConf, p.DirDriver} {
		if err := os.Rename(d, fmt.Sprintf("%s_old", d)); err != nil {
			p.tracker.Fail("backup", err)
			return err
		}

		p.tracker.messages <- fmt.Sprintf("%[1]s -> %[1]s_old", d)
		p.tracker.Accumulate("backup", 1)
	}
	return nil
}

func (p Porter) cleanup(restore bool) error {
	p.tracker.Start("cleanup", 2)

	if restore {
		p.tracker.messages <- "Restoring backups..."

		for _, d := range []string{p.DirConf, p.DirDriver} {
			if err := os.RemoveAll(d); err != nil {
				p.tracker.messages <- err.Error()
				p.tracker.Fail("cleanup", err)
				return err
			}

			if err := os.Rename(fmt.Sprintf("%s_old", d), d); err != nil {
				p.tracker.messages <- err.Error()
				p.tracker.Fail("cleanup", err)
				return err
			}

			p.tracker.messages <- fmt.Sprintf("%[1]s_old -> %[1]s", d)
			p.tracker.Accumulate("cleanup", 1)
		}
		return nil
	} else {
		p.tracker.messages <- "Cleaning up backups..."

		for _, d := range []string{p.DirConf, p.DirDriver} {
			p.tracker.messages <- fmt.Sprintf("Removing: %s_old", d)
			if err := os.RemoveAll(fmt.Sprintf("%s_old", d)); err != nil {
				// not able to removing backup is not a critical problem
				p.tracker.messages <- err.Error()
			} else {
				p.tracker.Accumulate("cleanup", 1)
			}
		}
		return nil
	}
}
