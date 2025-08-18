package status

type Status string

const (
	Pending   Status = "pending"
	Running   Status = "running"
	Completed Status = "completed"
	Failed    Status = "failed"
	Aborting  Status = "aborting"
	Aborted   Status = "aborted"
	Skiped    Status = "skiped"
	Speeded   Status = "speeded"
	Errored   Status = "errored"
)
