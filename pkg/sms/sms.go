package sms

import "context"

type Sms interface {
	Send(ctx context.Context, mobile, content string) error
}
