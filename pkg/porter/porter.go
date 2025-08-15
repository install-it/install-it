package porter

import (
	"context"
	"errors"
	"os"
	"path/filepath"
)

// Porter manages the porting process including export, import, and progress tracking.
type Porter struct {
	DirRoot    string             // Root directory for import/export operations
	Targets    []string           // Target directories to be backed up or compressed
	Message    chan string        // Channel for progress messages
	progresses []*Progress        // Slice of progress trackers for each step
	cancelFunc context.CancelFunc // Function to cancel ongoing operations
}

// Returns the current status of the porting process.
func (p Porter) Status() string {
	if len(p.progresses) == 0 {
		return "pending"
	}

	if some(p.progresses, func(p *Progress) bool { return p.Error == context.Canceled }) {
		if some(p.progresses, func(p *Progress) bool {
			return p.Status == "pending" || p.Status == "running"
		}) {
			return "aborting"
		}
		return "aborted"
	}

	if all(p.progresses, func(p *Progress) bool { return p.Status == "pending" }) {
		return "pending"
	}
	if all(p.progresses, func(p *Progress) bool { return p.Status == "completed" }) {
		return "completed"
	}
	if all(p.progresses, func(p *Progress) bool { return p.Status != "failed" }) {
		return "running"
	}
	return "failed"
}

// Cancels the ongoing porting process.
func (p Porter) Abort() error {
	if len(p.progresses) == 0 {
		return errors.New("porter: no started porting job")
	}

	switch p.Status() {
	case "aborting":
		return nil
	case "aborted":
		return errors.New("porter: already aborted")
	case "running":
		p.Message <- "Cancelling..."
		p.cancelFunc()
		return nil
	default:
		return errors.New("porter: no running porting job")
	}
}

// Returns the current progress and messages of the porting process.
func (p Porter) Progress() (Progresses, error) {
	if len(p.progresses) == 0 {
		return Progresses{}, errors.New("porter: no started porting job")
	}

	messageses := make([]string, 0, len(p.Message))
	for range len(p.Message) {
		messageses = append(messageses, <-p.Message)
	}

	progresses := make([]Progress, len(p.progresses))
	for i, prog := range p.progresses {
		progresses[i] = *prog
	}

	return Progresses{
		Progresses: progresses,
		Messages:   messageses,
		Status:     p.Status(),
		Error:      "",
	}, nil
}

// Compresses the target directories into a ZIP file at the destination.
func (p *Porter) Export(dest string) (err error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	p.cancelFunc = cancelFunc

	p.progresses = []*Progress{
		{context: ctx, message: p.Message, Name: "initialisation", Status: "pending"},
		{context: ctx, message: p.Message, Name: "compression", Status: "pending"},
	}
	defer p.exit()

	paths, err := func(tracker *Progress) (paths []string, err error) {
		tracker.Start(int64(len(p.Targets)))
		defer updateProgress(tracker, err)

		if cwd, err := os.Getwd(); err != nil {
			return nil, err
		} else {
			relpaths := []string{}
			for _, dir := range p.Targets {
				if rel, err := filepath.Rel(cwd, dir); err != nil {
					return relpaths, err
				} else {
					tracker.Accumulate(1)
					relpaths = append(relpaths, rel)
				}
			}
			return relpaths, nil

		}
	}(p.progresses[0])

	if err != nil {
		return err
	}
	return toZip(p.progresses[1], dest, paths...)
}

// Restores data from a ZIP file and cleans up or restores backups.
func (p *Porter) ImportFromFile(orig string) error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	p.cancelFunc = cancelFunc

	p.progresses = []*Progress{
		{context: ctx, message: p.Message, Name: "backup", Status: "pending"},
		{context: ctx, message: p.Message, Name: "decompression", Status: "pending"},
		{context: ctx, message: p.Message, Name: "cleanup", Status: "pending"},
	}
	defer p.exit()

	if err := backup(p.progresses[0], p.Targets); err != nil {
		return err
	}

	err := fromZip(p.progresses[1], orig, p.DirRoot)
	return errors.Join(err, cleanup(p.progresses[2], p.Targets, err != nil))
}

// Downloads a ZIP file from a URL and imports its contents.
func (p *Porter) ImportFromURL(url string) error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	p.cancelFunc = cancelFunc

	p.progresses = []*Progress{
		{context: ctx, message: p.Message, Name: "backup", Status: "pending"},
		{context: ctx, message: p.Message, Name: "download", Status: "pending"},
		{context: ctx, message: p.Message, Name: "decompression", Status: "pending"},
		{context: ctx, message: p.Message, Name: "cleanup", Status: "pending"},
	}
	defer p.exit()

	if err := backup(p.progresses[0], p.Targets); err != nil {
		return err
	}

	filename, err := download(p.progresses[1], url)
	if err != nil {
		return errors.Join(err, cleanup(p.progresses[3], p.Targets, true))
	}

	err = fromZip(p.progresses[2], filename, p.DirRoot)
	return errors.Join(err, cleanup(p.progresses[3], p.Targets, err != nil))
}

// Marks all pending progress steps as skipped.
func (p *Porter) exit() {
	for _, prog := range p.progresses {
		if prog.Status == "pending" {
			prog.Status = "skipped"
		}
	}
}
