package porter

import (
	"context"
	"errors"
	"os"
	"path/filepath"
)

type Porter struct {
	DirRoot    string
	Targets    []string
	Message    chan string
	progresses []*Progress
	cancelFunc context.CancelFunc
}

func (p Porter) Status() string {
	if len(p.progresses) == 0 {
		return "pending"
	}

	if some(p.progresses, func(p *Progress) bool { return p.Error == context.Canceled }) {
		if some(p.progresses, func(p *Progress) bool {
			return p.Status == "pending" || p.Status == "running"
		}) {
			return "aborting"
		} else {
			return "aborted"
		}
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

func (p Porter) Abort() error {
	if len(p.progresses) == 0 {
		return errors.New("porter: no started porting job")
	}

	if p.Status() == "aborting" {
		return nil
	}

	if p.Status() != "running" {
		return errors.New("porter: no running porting job")
	}

	if p.Status() == "aborted" {
		return errors.New("porter: already aborted")
	}

	p.Message <- "Cancelling..."
	p.cancelFunc()
	return nil
}

func (p Porter) Progress() (Progresses, error) {
	if len(p.progresses) == 0 {
		return Progresses{}, errors.New("porter: no started porting job")
	}

	messageses := make([]string, 0, len(p.Message))
	for range len(p.Message) {
		messageses = append(messageses, <-p.Message)
	}

	progresses := make([]Progress, 0, len(p.progresses))
	for _, prog := range p.progresses {
		progresses = append(progresses, *prog)
	}

	return Progresses{
		Progresses: progresses,
		Messages:   messageses,
		Status:     p.Status(),
		Error:      "",
	}, nil
}

func (p *Porter) Export(dest string) error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	p.cancelFunc = cancelFunc

	p.progresses = []*Progress{
		{context: ctx, message: p.Message, Name: "initialisation", Status: "pending"},
		{context: ctx, message: p.Message, Name: "compression", Status: "pending"},
	}
	defer p.exit()

	cwd, paths, err := func(tracker *Progress) (cwd string, paths []string, err error) {
		tracker.Start(1)
		defer updateProgress(tracker, err)

		if cwd, err := os.Getwd(); err != nil {
			return "", []string{}, err
		} else {
			if pathExe, err := os.Executable(); err != nil {
				return "", []string{}, err
			} else {
				root := filepath.Dir(pathExe)
				if cwd != root {
					os.Chdir(root)
					// defer func() { os.Chdir(cwd) }()
				}

				relpaths := []string{}
				for _, dir := range p.Targets {
					if rel, err := filepath.Rel(root, dir); err != nil {
						return cwd, []string{}, err
					} else {
						relpaths = append(relpaths, rel)
					}
				}
				return cwd, relpaths, nil
			}
		}
	}(p.progresses[0])

	defer func() {
		if cwd != "" {
			os.Chdir(cwd)
		}
	}()

	if err != nil {
		return err
	}

	return toZip(p.progresses[1], dest, paths...)

}

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

	if filename, err := download(p.progresses[1], url); err != nil {
		return errors.Join(err, cleanup(p.progresses[3], p.Targets, true))
	} else {
		err = fromZip(p.progresses[2], filename, p.DirRoot)
		return errors.Join(err, cleanup(p.progresses[3], p.Targets, err != nil))
	}
}

func (p *Porter) exit() {
	for _, p := range p.progresses {
		if p.Status == "pending" {
			p.Status = "skiped"
		}
	}
}
