package report

import "context"


type Report struct {
	Location string
	Format string
}

type Generator interface {
	Generate(ctx context.Context, reportID string) (*Report, error)
}