package executor

import (
	"context"

	"github.com/martin0b101/ilmarinen/internal/job"
)


type Executor interface{
	Execute(ctx context.Context, job *job.Job)
}