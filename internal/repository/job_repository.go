package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/martin0b101/ilmarinen/internal/model"
)



type JobRepository interface{
	Create(ctx context.Context, job model.JobModel) error
	Update(ctx context.Context, id uuid.UUID, job model.JobModel) error
	ClaimNext(ctx context.Context) (*model.JobModel, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.JobModel, error)
}