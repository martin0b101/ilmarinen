package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type JobModel struct{
    ID          uuid.UUID       `db:"id"`
    Type        string          `db:"job_type"`
    Payload     json.RawMessage `db:"payload"`
    Status      string          `db:"status"`
    Attempts    int             `db:"attempts"`
    MaxAttempts int             `db:"max_attempts"`
    ScheduledAt time.Time       `db:"scheduled_at"`
    CreatedAt   time.Time       `db:"created_at"`
    UpdatedAt   time.Time       `db:"updated_at"`
    CompletedAt *time.Time      `db:"completed_at"`
}