package email

import (
	"context"
	"fmt"
)

var _ Sender = (*SmtpSender)(nil)

type SmtpSender struct {
}

// Send implements [Sender].
func (s *SmtpSender) Send(ctx context.Context, to string, subject string, body string) error {

	fmt.Printf("Sending msg to: %s\n", to)
	return nil
}
