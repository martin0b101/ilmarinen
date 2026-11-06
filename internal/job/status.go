package job

const (
	Pending   JobStatus = "PENDING"
	Running   JobStatus = "RUNNING"
	Completed JobStatus = "COMPLETED"
	Canceled  JobStatus = "CANCELED"
	Failed    JobStatus = "FAILED"
)

type JobStatus string

var statusTransition map[JobStatus][]JobStatus = map[JobStatus][]JobStatus{
	Pending:   {Running, Canceled},
	Running:   {Completed, Pending, Failed, Canceled},
	Failed:    {Pending},
	Canceled:  {Pending},
	Completed: {},
}
