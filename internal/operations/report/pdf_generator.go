package report

import (
	"context"
	"fmt"
)

var _ Generator = (*PdfGenerator)(nil)

type PdfGenerator struct {
}

// Generate implements [Generator].
func (p *PdfGenerator) Generate(ctx context.Context, reportID string) (*Report, error) {
	
	fmt.Printf("Generating %s report....\n", reportID)
	return &Report{
		Location: "path-to-report",
		Format: ".pdf",
	},nil
}

