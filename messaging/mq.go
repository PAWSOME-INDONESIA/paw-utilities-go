package messaging

import (
	"context"
)

type CallbackFunc func([]byte) error

type Queue interface {
	ReadWithContext(context.Context, string, []CallbackFunc) error
	Read(string, []CallbackFunc) error
	PublishWithContext(context.Context, string, string) error
	Publish(string, string) error
	Close() error
}
