package porter

import (
	"context"
	"time"
)

// Progress counts the number of bytes written to it.
// It implements to the io.Writer interface and we can pass this into io.TeeReader() which will report Progress on each write cycle.
type Progress struct {
	Status  string    `json:"status"`
	Total   int64     `json:"total"`
	Current int64     `json:"current"`
	StartAt time.Time `json:"startAt"`
	Error   error     `json:"error"`
}

func (p *Progress) Write(b []byte) (int, error) {
	n := len(b)
	p.Current += int64(n)
	return n, nil
}

func (p *Progress) Start(total int64) {
	p.StartAt = time.Now()
	p.Status = "running"
	p.Total = total
}

func (p *Progress) Accumulate(current int64) {
	p.Current += current
	if p.Total == p.Current {
		p.Status = "completed"
	}
}

func (p *Progress) Complete() {
	p.Current = p.Total
	p.Status = "completed"
}

func (p *Progress) Fail(err error) {
	p.Error = err

	if err == context.Canceled {
		p.Status = "aborted"
	} else {
		p.Status = "failed"
	}
}

// Type binding to the frontend progress query
type Progresses struct {
	Progresses []Progress `json:"tasks"`
	Messages   []string   `json:"message"`
	Status     string     `json:"status"`
	Error      string     `json:"error"`
}

type ProgressTracker struct {
	progs    map[string]*Progress
	messages chan string
	err      error
	ctx      context.Context
}

func (pt ProgressTracker) Get(name string) *Progress {
	return pt.progs[name]
}

func (pt *ProgressTracker) Add(name string, tracker *Progress) {
	pt.progs[name] = tracker
}

func (pt *ProgressTracker) Start(name string, total int64) {
	pt.progs[name].StartAt = time.Now()
	pt.progs[name].Status = "running"
	pt.progs[name].Total = total
}

func (pt *ProgressTracker) Accumulate(name string, current int64) {
	pt.progs[name].Current += current
	if pt.progs[name].Total == pt.progs[name].Current {
		pt.progs[name].Status = "completed"
	}
}

func (pt *ProgressTracker) Complete(name string) {
	pt.progs[name].Current = pt.progs[name].Total
	pt.progs[name].Status = "completed"
}

func (pt *ProgressTracker) Fail(name string, err error) {
	pt.err = err

	if err == context.Canceled {
		pt.progs[name].Status = "aborted"
	} else {
		pt.progs[name].Status = "failed"
	}
}

func (pt *ProgressTracker) Exit() {
	for _, tracker := range pt.progs {
		if tracker.Status == "pending" {
			tracker.Status = "skiped"
		}
	}
}

func NewProgressTracker(ctx context.Context, trackers map[string]*Progress) *ProgressTracker {
	return &ProgressTracker{
		progs:    trackers,
		ctx:      ctx,
		messages: make(chan string, 128),
	}
}
