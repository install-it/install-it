package porter

import (
	"context"
	"install-it/pkg/status"
	"sync"
	"time"
)

// JobSnapshot is a point-in-time view of the current job, polled by the frontend.
type JobSnapshot struct {
	Status   status.Status `json:"status"`   // pending|running|completed|failed|aborted
	Step     string        `json:"step"`      // "download"|"backup"|"extract"|"cleanup"|"" when idle
	Progress float64       `json:"progress"`  // 0.0 to 1.0
	Messages []string      `json:"messages"`  // recent messages (tail)
}

// job tracks the state of a single export/import operation.
type job struct {
	mu       sync.Mutex
	status   status.Status
	step     string
	progress float64
	messages []string // ring buffer, keep last 100
	startAt  time.Time

	ctx       context.Context
	cancel    context.CancelFunc
	messageCh chan string // buffered channel for receiving messages from worker
}

func newJob() *job {
	ctx, cancel := context.WithCancel(context.Background())
	return &job{
		status:    status.Pending,
		ctx:       ctx,
		cancel:    cancel,
		messageCh: make(chan string, 512),
		messages:  make([]string, 0, 100),
	}
}

func (j *job) start() {
	j.mu.Lock()
	j.status = status.Running
	j.startAt = time.Now()
	j.mu.Unlock()
}

func (j *job) setStep(name string) {
	j.mu.Lock()
	j.step = name
	j.mu.Unlock()
}

func (j *job) setProgress(p float64) {
	j.mu.Lock()
	if p < 0 {
		p = 0
	}
	if p > 1 {
		p = 1
	}
	j.progress = p
	j.mu.Unlock()
}

// msg sends a message to the channel (non-blocking via select with default).
func (j *job) msg(s string) {
	select {
	case j.messageCh <- s:
	default:
	}
}

func (j *job) complete() {
	j.mu.Lock()
	j.status = status.Completed
	j.progress = 1.0
	j.mu.Unlock()
}

// fail sets status=Failed (or Aborted if err is context.Canceled).
func (j *job) fail(err error) {
	j.mu.Lock()
	if err == context.Canceled {
		j.status = status.Aborted
	} else {
		j.status = status.Failed
	}
	j.mu.Unlock()
}

// snapshot drains messageCh into messages (keep last 100) and returns a snapshot.
func (j *job) snapshot() JobSnapshot {
	for {
		select {
		case m := <-j.messageCh:
			j.messages = append(j.messages, m)
		default:
			goto done
		}
	}
done:

	if len(j.messages) > 100 {
		j.messages = j.messages[len(j.messages)-100:]
	}

	j.mu.Lock()
	defer j.mu.Unlock()

	msgs := make([]string, len(j.messages))
	copy(msgs, j.messages)

	return JobSnapshot{
		Status:   j.status,
		Step:     j.step,
		Progress: j.progress,
		Messages: msgs,
	}
}


