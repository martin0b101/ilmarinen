package executor

import (
	"context"

	pkgJob "github.com/martin0b101/ilmarinen/internal/job"
	"github.com/martin0b101/ilmarinen/internal/operations/email"
)

var _ Executor = (*EmailExecutor)(nil)

type EmailExecutor struct{
	Sender email.SmtpSender
}

// Execute implements [Executor].
func (e *EmailExecutor) Execute(ctx context.Context, job *pkgJob.Job) {
	// unrefrence payload
	payload, ok := job.Payload().(*pkgJob.EmailPayload)
	if !ok {
		return
	}
	e.Sender.Send(ctx, payload.To, payload.Subject, payload.Body)
}

