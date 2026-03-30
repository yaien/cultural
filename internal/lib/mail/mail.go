package mail

import "context"

type Mail interface {
	Send(ctx context.Context, email *Email) error
}

type Recipient struct {
	Name  string
	Email string
}

type Email struct {
	To       Recipient
	From     Recipient
	Subject  string
	Body     string
	Category string
}
