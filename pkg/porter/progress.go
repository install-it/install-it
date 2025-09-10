package porter

import (
	"context"
	"install-it/pkg/status"
	"time"
)

// Progress tracks the status and metrics of a single task.
// It implements io.Writer to allow tracking of byte-based progress (e.g., during downloads).
type Progress struct {
	Name    string          `json:"name"`   // Name of the task
	Status  status.Status   `json:"status"` // Current status: "pending", "running", "completed", "failed", "aborted"
	Total   int64           `json:"total"`
	Current int64           `json:"current"`
	StartAt time.Time       `json:"startAt"` // Timestamp when the task started
	Error   error           `json:"error"`   // Error encountered during execution, if any
	message chan string     // Channel for sending progress messages
	context context.Context // Context for cancellation and timeout control
}

// Implements the io.Writer interface.
// It updates the Current progress based on the number of bytes written.
func (p *Progress) Write(b []byte) (int, error) {
	n := len(b)
	p.Current += int64(n)
	return n, nil
}

// Initializes the progress tracking with a total byte count.
func (p *Progress) Start(total int64) {
	p.StartAt = time.Now()
	p.Status = status.Running
	p.Total = total
}

// Adds a given number of bytes to the current progress.
func (p *Progress) Accumulate(current int64) {
	p.Current += current
}

// Marks the progress as completed and sets Current to Total.
func (p *Progress) Complete() {
	p.Current = p.Total
	p.Status = status.Completed
}

// Marks the progress as failed or aborted depending on the error type.
func (p *Progress) Fail(err error) {
	p.Error = err
	if err == context.Canceled {
		p.Status = status.Aborted
	} else {
		p.Status = status.Failed
	}
}

// Type binding to the frontend progress query
type Progresses struct {
	Progresses []Progress    `json:"tasks"`    // List of individual task progress
	Messages   []string      `json:"messages"` // Collected messages from all tasks
	Status     status.Status `json:"status"`   // Overall status of the porting process
}
