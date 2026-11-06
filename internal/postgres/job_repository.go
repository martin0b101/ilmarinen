package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/martin0b101/ilmarinen/internal/model"
	"github.com/martin0b101/ilmarinen/internal/repository"
	pkgErrors "github.com/martin0b101/ilmarinen/internal/errors"

)

type JobRepository struct {
	db *sql.DB
}

// ClaimNext implements [repository.JobRepository].
func (jr *JobRepository) ClaimNext(ctx context.Context) (*model.JobModel, error) {
    query := `
        WITH next_job AS (
            SELECT id
            FROM jobs
            WHERE status = 'PENDING'
              AND scheduled_at <= NOW()
            ORDER BY scheduled_at ASC, created_at ASC
            LIMIT 1
            FOR UPDATE SKIP LOCKED
        )
        UPDATE jobs AS j
        SET
            status = 'RUNNING',
            updated_at = NOW()
        FROM next_job
        WHERE j.id = next_job.id
        RETURNING
            j.id,
            j.job_type,
            j.payload,
            j.status,
            j.attempts,
            j.max_attempts,
            j.scheduled_at,
            j.created_at,
            j.updated_at,
            j.completed_at
    `

    var job model.JobModel

    err := jr.db.QueryRowContext(ctx, query).Scan(
        &job.ID,
        &job.Type,
        &job.Payload,
        &job.Status,
        &job.Attempts,
        &job.MaxAttempts,
        &job.ScheduledAt,
        &job.CreatedAt,
        &job.UpdatedAt,
        &job.CompletedAt,
    )

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, pkgErrors.ErrNoJobsAvailable
        }

        return nil, fmt.Errorf("claim next job: %w", err)
    }

    return &job, nil
}

// Create implements [repository.JobRepository].
func (j *JobRepository) Create(ctx context.Context, job model.JobModel) error {
	query := `
        INSERT INTO jobs (
            id,
            job_type,
            payload,
            status,
            attempts,
            max_attempts,
            scheduled_at,
            created_at,
            updated_at,
            completed_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `

	_, err := j.db.ExecContext(
		ctx,
		query,
		job.ID,
		job.Type,
		job.Payload,
		job.Status,
		job.Attempts,
		job.MaxAttempts,
		job.ScheduledAt,
		job.CreatedAt,
		job.UpdatedAt,
		job.CompletedAt,
	)

	if err != nil {
		return fmt.Errorf("create job: %w", err)
	}

	return nil
}

// GetByID implements [repository.JobRepository].
func (j *JobRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.JobModel, error) {
	query := `
        SELECT
            id,
            job_type,
            payload,
            status,
            attempts,
            max_attempts,
            scheduled_at,
            created_at,
            updated_at,
            completed_at
        FROM jobs
        WHERE id = $1
    `

	var job model.JobModel

	err := j.db.QueryRowContext(ctx, query, id).Scan(
		&job.ID,
		&job.Type,
		&job.Payload,
		&job.Status,
		&job.Attempts,
		&job.MaxAttempts,
		&job.ScheduledAt,
		&job.CreatedAt,
		&job.UpdatedAt,
		&job.CompletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("job with ID %s not found: %w", id, err)
		}

		return nil, fmt.Errorf("get job by ID: %w", err)
	}

	return &job, nil
}

// Update implements [repository.JobRepository].
func (j *JobRepository) Update(ctx context.Context, id uuid.UUID, job model.JobModel) error {

	query := `
    UPDATE jobs
    SET
        status = $1,
        attempts = $2,
        updated_at = $3,
        completed_at = $4
    WHERE id = $5
`

	result, err := j.db.ExecContext(
		ctx,
		query,
		job.Status,
		job.Attempts,
		job.UpdatedAt,
		job.CompletedAt,
		job.ID,
	)
	if err != nil {
		return fmt.Errorf("update job: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("job with ID %s not found", job.ID)
	}

	return nil
}

var _ repository.JobRepository = (*JobRepository)(nil)

func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{
		db: db,
	}
}
