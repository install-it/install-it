package porter

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Porter struct {
	DirRoot    string
	DirConf    string
	DirDriver  string
	tracker    *ProgressTracker
	cancelFunc context.CancelFunc
}

func (p Porter) Status() string {
	if p.tracker == nil {
		return "pending"
	}

	if p.tracker.ctx.Err() == context.Canceled {
		if some(p.tracker.progresses, func(p *Progress) bool {
			return p.Status == "pending" || p.Status == "running"
		}) {
			return "aborting"
		} else {
			return "aborted"
		}
	} else if all(p.tracker.progresses, func(p *Progress) bool { return p.Status == "pending" }) {
		return "pending"
	} else if all(p.tracker.progresses, func(p *Progress) bool { return p.Status == "completed" }) {
		return "completed"
	} else if all(p.tracker.progresses, func(p *Progress) bool { return p.Status != "failed" }) {
		// "aborting" and "aborted" is eliminated in the above conditions
		return "running"
	} else {
		return "failed"
	}
}

func (p Porter) Abort() error {
	if p.tracker == nil {
		return errors.New("porter: no started porting job")
	}

	if p.Status() == "aborting" {
		return nil
	}

	if p.Status() != "running" {
		return errors.New("porter: no running porting job")
	}

	if p.tracker.ctx.Err() == context.Canceled {
		return errors.New("porter: already aborted")
	}

	p.tracker.messages <- "Cancelling..."

	p.cancelFunc()

	return nil
}

func (p Porter) Progress() (Progresses, error) {
	if p.tracker == nil {
		return Progresses{}, errors.New("porter: no started porting job")
	}

	messages := make([]string, len(p.tracker.messages))
	for range len(p.tracker.messages) {
		messages = append(messages, <-p.tracker.messages)
	}

	tasks := make([]Progress, 0)
	for _, t := range p.tracker.progresses {
		tasks = append(tasks, *t)
	}

	var error_str string
	if p.tracker.err != nil {
		error_str = p.tracker.err.Error()
	}

	return Progresses{
		Progresses: tasks,
		Messages:   messages,
		Status:     p.Status(),
		Error:      error_str,
	}, nil
}

func (p *Porter) Export(dest string) error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	p.cancelFunc = cancelFunc

	p.tracker = NewProgressTracker(ctx, []*Progress{
		{Name: "initialisation", Status: "pending"},
		{Name: "compression", Status: "pending"},
	})
	defer p.tracker.SkipAllPendings()

	p.tracker.Start("initialisation", 1)

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if pathExe, err := os.Executable(); err != nil {
		return err
	} else {
		if cwd != filepath.Dir(pathExe) {
			os.Chdir(filepath.Dir(pathExe))
			defer os.Chdir(cwd)

			cwd = filepath.Dir(pathExe)
		}
	}

	relDirConf, err := filepath.Rel(cwd, p.DirConf)
	if err != nil {
		return err
	}

	relDirDir, err := filepath.Rel(cwd, p.DirDriver)
	if err != nil {
		return err
	}

	p.tracker.Complete("initialisation")

	return toZip(p.tracker, dest, relDirConf, relDirDir)

}

func (p *Porter) ImportFromFile(orig string, igoreSetting bool) error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	p.cancelFunc = cancelFunc

	p.tracker = NewProgressTracker(ctx, []*Progress{
		{Name: "backup", Status: "pending"},
		{Name: "decompression", Status: "pending"},
		{Name: "cleanup", Status: "pending"},
	})
	defer p.tracker.SkipAllPendings()

	if err := p.backup(); err != nil {
		return err
	}

	err := fromZip(p.tracker, orig, p.DirRoot)
	if err == nil && igoreSetting {
		os.Rename(
			filepath.Join(p.DirConf+"_old", "setting.json"),
			filepath.Join(p.DirConf, "setting.json"))
	}

	return errors.Join(err, p.cleanup(err != nil))
}

func (p *Porter) ImportFromURL(url string, igoreSetting bool) error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	p.cancelFunc = cancelFunc

	p.tracker = NewProgressTracker(ctx, []*Progress{
		{Name: "initialisation", Status: "pending"},
		{Name: "backup", Status: "pending"},
		{Name: "download", Status: "pending"},
		{Name: "decompression", Status: "pending"},
		{Name: "cleanup", Status: "pending"},
	})
	defer p.tracker.SkipAllPendings()

	p.tracker.Start("initialisation", 1)

	var filename string
	for {
		if filename != "" {
			if _, err := os.Stat(filename); err != nil {
				break
			}
		}

		rb := make([]byte, 4)
		rand.Read(rb)
		filename = fmt.Sprintf("%s%s.zip", os.TempDir(), hex.EncodeToString(rb))
	}

	p.tracker.Complete("initialisation")

	if err := p.backup(); err != nil {
		return err
	}

	if err := download(p.tracker, url, filename); err != nil {
		return errors.Join(err, p.cleanup(true))
	} else {
		err = fromZip(p.tracker, filename, p.DirRoot)
		return errors.Join(err, p.cleanup(err != nil))
	}
}
