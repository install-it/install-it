package porter

import (
	"context"
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
	Error   error     `json:"error"`
	message chan string
	context context.Context
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
