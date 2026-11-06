package job

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testEmailMsg = json.RawMessage(`{
    "to": "test@example.com",
    "subject": "Hello",
    "body": "Test email"
}`)

var defaultMaxAttempts int = 3
var defaultRetyrDelay time.Duration = 5 * time.Second

func SetupEmailJob() *Job {
	return NewJob(EmailJobType, testEmailMsg, defaultMaxAttempts, defaultRetyrDelay)
}

func TestTransitions(t *testing.T) {

	tests := []struct {
		name             string
		setup            func() *Job
		action           func(j *Job) error
		expectedStatus   JobStatus
		expectedAttempts *int
		validate         func(j *Job) bool
		wantErr          bool
		expectedErr      string
	}{
		{
			name:  "change transitions from pending to running",
			setup: SetupEmailJob,
			action: func(j *Job) error {
				return j.Start()
			},
			expectedStatus: Running,
			wantErr:        false,
		},
		{
			name: "not valid transitions",
			setup: func() *Job {
				job := NewJob(EmailJobType, testEmailMsg, defaultMaxAttempts, defaultRetyrDelay)
				_ = job.Start()
				return job
			},
			action: func(j *Job) error {
				return j.Start()
			},
			expectedStatus: Running,
			wantErr:        true,
			expectedErr:    fmt.Sprintf("transition from %s to %s is not valid", Running, Running),
		},
		{
			name:  "starting a job increments attempts",
			setup: SetupEmailJob,
			action: func(j *Job) error {
				return j.Start()
			},
			expectedAttempts: new(1),
			expectedStatus:   Running,
			wantErr:          false,
		},
		{
			name:  "too much attempts",
			setup: SetupEmailJob,
			action: func(j *Job) error {
				// 1. attempt
				_ = j.Start()
				_ = j.Fail()
				_ = j.Retry()
				// 3. attempt
				_ = j.Start()
				_ = j.Fail()
				_ = j.Retry()
				// 3. attempt
				_ = j.Start()
				_ = j.Fail()
				_ = j.Retry()
				// failed on 4 attempt
				return j.Start()
			},
			wantErr:     true,
			expectedErr: fmt.Sprintf("too many attempts to start"),
		},
		{
			name:  "complete transiton updates completedAt",
			setup: SetupEmailJob,
			action: func(j *Job) error {
				_ = j.Start()
				return j.Complete()
			},
			expectedStatus: Completed,
			validate: func(j *Job) bool {
				return j.completedAt != nil
			},
			wantErr: false,
		},
		{
			name:  "canceled pending job",
			setup: SetupEmailJob,
			action: func(j *Job) error {
				return j.Cancel()
			},
			expectedStatus: Canceled,
			wantErr:        false,
		},
		{
			name:  "cancel running job",
			setup: SetupEmailJob,
			action: func(j *Job) error {
				_ = j.Start()
				return j.Cancel()
			},
			expectedStatus: Canceled,
			wantErr:        false,
		},
		{
			name:  "cancel failed job - can not be canceled",
			setup: SetupEmailJob,
			action: func(j *Job) error {
				_ = j.Start()
				_ = j.Fail()
				fmt.Printf("JOB STATUS: %s\n", j.status)
				return j.Cancel()
			},
			wantErr:     true,
			expectedErr: fmt.Sprintf("transition from %s to %s is not valid", Failed, Canceled),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			job := tt.setup()
			err := tt.action(job)

			if tt.wantErr {
				if tt.expectedStatus != "" {
					assert.Equal(t, tt.expectedStatus, job.Status())
				}
				assert.EqualError(t, err, tt.expectedErr)
				return
			}
			if tt.expectedAttempts != nil {
				assert.Equal(t, *tt.expectedAttempts, job.attempts)
			}

			if tt.validate != nil {
				assert.True(t, tt.validate(job))
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, job.Status())
		})
	}
}

func TestStartUpdatesUpdatedAt(t *testing.T) {
	job := SetupEmailJob()
	previousUpdatedAt := job.updatedAt

	err := job.Start()

	assert.NoError(t, err)
	assert.True(t, job.updatedAt.After(previousUpdatedAt))
}

func TestRetryAdds5sToScheduledAt(t *testing.T) {
	job := SetupEmailJob()

	_ = job.Start()
	_ = job.Fail()

	retryTime := time.Now()
	err := job.Retry()

	assert.NoError(t, err)
	assert.InDelta(t,
		retryTime.Add(5*time.Second).UnixMilli(),
		job.scheduledAt.UnixMilli(),
		100)
}
