package porter_new

import (
	"context"
	"fmt"
	"time"
)

// Progress counts the number of bytes written to it.
// It implements to the io.Writer interface and we can pass this into io.TeeReader() which will report Progress on each write cycle.
type Progress struct {
	Name    string    `json:"name"`
	Status  string    `json:"status"`
	Total   int64     `json:"total"`
	Current int64     `json:"current"`
	StartAt time.Time `json:"startAt"`
}

func (p *Progress) Write(b []byte) (int, error) {
	n := len(b)
	p.Current += int64(n)
	return n, nil
}

// Type binding to the frontend progress query
type Progresses struct {
	Progresses []Progress `json:"tasks"`
	Messages   []string   `json:"message"`
	Status     string     `json:"status"`
	Error      string     `json:"error"`
}

type ProgressTracker struct {
	progresses []*Progress
	messages   chan string
	err        error
	ctx        context.Context
}

func (pt ProgressTracker) indexOf(name string) (int, error) {
	for index, progress := range pt.progresses {
		if progress.Name == name {
			return index, nil
		}
	}
	return -1, fmt.Errorf("porter: cannot found progress with %s", name)
}

func (pt ProgressTracker) Get(name string) (*Progress, error) {
	if index, err := pt.indexOf(name); err != nil {
		return nil, err
	} else {
		return pt.progresses[index], nil
	}
}

func (pt *ProgressTracker) Start(name string, total int64) {
	if index, err := pt.indexOf(name); err == nil {
		pt.progresses[index].Total = total
		pt.progresses[index].Status = "running"
		pt.progresses[index].StartAt = time.Now()
	} else {
		panic(err)
	}
}

func (pt *ProgressTracker) Accumulate(name string, current int64) {
	if index, err := pt.indexOf(name); err == nil {
		pt.progresses[index].Current += current

		if pt.progresses[index].Total == pt.progresses[index].Current {
			pt.progresses[index].Status = "completed"
		}
	} else {
		panic(err)
	}
}

func (pt *ProgressTracker) Complete(name string) {
	if index, err := pt.indexOf(name); err == nil {
		pt.progresses[index].Current = pt.progresses[index].Total
		pt.progresses[index].Status = "completed"
	} else {
		panic(err)
	}
}

func (pt *ProgressTracker) Fail(name string, err error) {
	pt.err = err

	if index, err := pt.indexOf(name); err == nil {
		if err == context.Canceled {
			pt.progresses[index].Status = "aborted"
		} else {
			pt.progresses[index].Status = "failed"
		}
	} else {
		panic(err)
	}
}

func (pt *ProgressTracker) SkipAllPendings() {
	for _, progress := range pt.progresses {
		if progress.Status == "pending" {
			progress.Status = "skiped"
		}
	}
}

func NewProgressTracker(ctx context.Context, progresses []*Progress) *ProgressTracker {
	return &ProgressTracker{
		progresses: progresses,
		ctx:        ctx,
		messages:   make(chan string, 128),
	}
}
