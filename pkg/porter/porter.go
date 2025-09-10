package porter

import (
	"context"
	"errors"
	"install-it/pkg/status"
	"install-it/pkg/utils"
	"os"
	"path/filepath"
	"strings"
)

// Porter manages the porting process including export, import, and progress tracking.
type Porter struct {
	DirRoot string   // Root directory for import/export operations
	Targets []string // Target directories to be backed up or compressed

	Message    chan string // Channel for progress messages
	progresses []*Progress // Slice of progress trackers for each step

	ctx        context.Context
	cancelFunc context.CancelFunc // Function to cancel ongoing operations
}

// Returns the current status of the porting process.
func (p Porter) Status() status.Status {
	if len(p.progresses) == 0 {
		return status.Pending
	}

	if p.ctx.Err() == context.Canceled {
		if utils.Some(p.progresses, func(p *Progress) bool {
			return strings.Contains(string(p.Status), "ing")
		}) {
			return status.Aborting
		}
		return status.Aborted
	}

	if utils.All(p.progresses, func(p *Progress) bool { return p.Status == status.Pending }) {
		return status.Pending
	}
	if utils.All(p.progresses, func(p *Progress) bool { return p.Status == status.Completed }) {
		return status.Completed
	}
	if utils.All(p.progresses, func(p *Progress) bool { return p.Status != status.Failed }) {
		return status.Running
	}
	return status.Failed
}

// Cancels the ongoing porting process.
func (p Porter) Abort() error {
	if len(p.progresses) == 0 || p.cancelFunc == nil {
		return errors.New("porter: no started porting job")
	}

	switch p.Status() {
	case status.Aborting:
		return nil
	case status.Aborted:
		return errors.New("porter: already aborted")
	case status.Running:
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
	}, nil
}

// Compresses the target directories into a ZIP file at the destination.
func (p *Porter) Export(dest string) (err error) {
	p.ctx, p.cancelFunc = context.WithCancel(context.Background())

	p.progresses = []*Progress{
		{context: p.ctx, message: p.Message, Name: "initialisation", Status: status.Pending},
		{context: p.ctx, message: p.Message, Name: "compression", Status: status.Pending},
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
	p.ctx, p.cancelFunc = context.WithCancel(context.Background())

	p.progresses = []*Progress{
		{context: p.ctx, message: p.Message, Name: "backup", Status: status.Pending},
		{context: p.ctx, message: p.Message, Name: "decompression", Status: status.Pending},
		{context: p.ctx, message: p.Message, Name: "cleanup", Status: status.Pending},
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
	p.ctx, p.cancelFunc = context.WithCancel(context.Background())

	p.progresses = []*Progress{
		{context: p.ctx, message: p.Message, Name: "backup", Status: status.Pending},
		{context: p.ctx, message: p.Message, Name: "download", Status: status.Pending},
		{context: p.ctx, message: p.Message, Name: "decompression", Status: status.Pending},
		{context: p.ctx, message: p.Message, Name: "cleanup", Status: status.Pending},
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
		if prog.Status == status.Pending {
			prog.Status = status.Skiped
		}
	}
}
