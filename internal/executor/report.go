package executor

import (
	"context"

	pkgJob "github.com/martin0b101/ilmarinen/internal/job"
	"github.com/martin0b101/ilmarinen/internal/operations/report"
)

var _ Executor = (*ReportExecutor)(nil)

type ReportExecutor struct{
	PdfGenerator report.PdfGenerator 
}

// Execute implements [Executor].
func (r *ReportExecutor) Execute(ctx context.Context, job *pkgJob.Job) {
	
	payload, ok := job.Payload().(*pkgJob.ReportPayload)
	if !ok {
		return
	}
	r.PdfGenerator.Generate(ctx, payload.ReportId)
}
