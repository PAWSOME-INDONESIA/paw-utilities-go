package messaging

import (
	"context"
)

type Queue interface {
	ReadWithContext(context.Context, string, []CallbackFunc) error
	Read(string, []CallbackFunc) error
	PublishWithContext(context.Context, string, string) error
	Publish(string, string) error
}
