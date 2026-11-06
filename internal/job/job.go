package job

import (
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
)

const (
	EmailJobType  JobType = "EMAIL"
	ReportJobType JobType = "REPORT"
)

type Payload interface {
	Validate() error
}

type JobType string

type Job struct {
	id          uuid.UUID
	jobType     JobType
	payload     Payload
	status      JobStatus
	attempts    int
	maxAttempts int
	scheduledAt time.Time
	createdAt   time.Time
	updatedAt   time.Time
	completedAt *time.Time
	retryDelay  time.Duration
}

type EmailPayload struct {
	To      string
	Subject string
	Body    string
}

func NewEmailPayload(payload json.RawMessage) (*EmailPayload, error) {
	if payload == nil {
		return nil, fmt.Errorf("paylod is empty")
	}
	var email EmailPayload
	err := json.Unmarshal(payload, &email)
	if err != nil {
		return nil, err
	}

	err = email.Validate()
	if err != nil {
		return nil, err
	}

	return &email, nil
}

func (ep *EmailPayload) Validate() error {
	return nil
}

type ReportPayload struct {
	ReportId string
}

func NewReportPayload(payload json.RawMessage) (*ReportPayload, error) {
	if payload == nil {
		return nil, fmt.Errorf("paylod is empty")
	}
	var report ReportPayload
	err := json.Unmarshal(payload, &report)
	if err != nil {
		return nil, err
	}

	err = report.Validate()
	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (rp *ReportPayload) Validate() error {
	return nil
}

func NewJob(jobType JobType, payload json.RawMessage, maxAttempts int, retryDelay time.Duration) *Job {

	job := &Job{
		id:          uuid.New(),
		jobType:     jobType,
		status:      Pending,
		attempts:    0,
		maxAttempts: maxAttempts,
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
		scheduledAt: time.Now(),
		retryDelay:  retryDelay,
	}

	switch jobType {
	case ReportJobType:
		payload, err := NewReportPayload(payload)
		if err != nil {
			return nil
		}
		job.payload = payload
		break
	case EmailJobType:
		payload, err := NewReportPayload(payload)
		if err != nil {
			return nil
		}
		job.payload = payload
		break
	default:
		return nil
	}

	return job
}

func (j *Job) Status() JobStatus {
	return j.status
}

func (j *Job) Id() uuid.UUID {
	return j.id
}

func (j *Job) Type() JobType {
	return j.jobType
}

func (j *Job) Payload() Payload {
	return j.payload
}

func (j *Job) Attempts() int {
	return j.attempts
}

func (j *Job) ScheduledAt() time.Time {
	return j.scheduledAt
}

func (j *Job) Validate() error {
	if j.id == uuid.Nil {
		return fmt.Errorf("id is empty")
	}

	if j.Type() != EmailJobType || j.Type() != ReportJobType {
		return fmt.Errorf("type is not supported")
	}

	// .. and so on
	return nil
}

func (j *Job) setUpdateAt() {
	j.updatedAt = time.Now()
}

func (j *Job) setCompletedAt() {
	now := time.Now()
	j.completedAt = &now
}

func (j *Job) setBackOffTime() {
	backoffTime := j.retryDelay * time.Duration(1<<(j.attempts-1))
	j.scheduledAt = time.Now().Add(backoffTime)
}

func (j *Job) setStatus(status JobStatus) error {
	err := j.validateTransition(status)
	if err != nil {
		return err
	}
	j.status = status
	j.setUpdateAt()
	if status == Completed {
		j.setCompletedAt()
	}
	return nil
}

func (j *Job) validateTransition(nextStatus JobStatus) error {
	validTranstions := statusTransition[j.status]
	if !slices.Contains(validTranstions, nextStatus) {
		return fmt.Errorf("transition from %s to %s is not valid", j.status, nextStatus)
	}
	return nil
}

func (j *Job) Complete() error {
	return j.setStatus(Completed)
}
func (j *Job) Cancel() error {
	return j.setStatus(Canceled)
}
func (j *Job) Fail() error {
	return j.setStatus(Failed)
}

func (j *Job) Start() error {
	if j.attempts >= j.maxAttempts {
		return fmt.Errorf("too many attempts to start")
	}
	err := j.setStatus(Running)
	if err != nil {
		return err
	}
	j.attempts += 1
	return nil
}

func (j *Job) Retry() error {
	if j.attempts == 0 {
		return fmt.Errorf("no execution was attempted")
	}
	if j.attempts >= j.maxAttempts {
		return fmt.Errorf("maximum job attempts reached")
	}
	err := j.setStatus(Pending)
	if err != nil {
		return err
	}
	j.setBackOffTime()
	return nil
}
