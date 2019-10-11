package messaging

import (
	"context"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/util"
)

type CallbackFunc func([]byte) error

type Queue interface {
	util.Ping
	ReadWithContext(context.Context, string, []CallbackFunc) error
	Read(string, []CallbackFunc) error
	PublishWithContext(context.Context, string, string) error
	Publish(string, string) error
	Close() error
}

type MessagingQueue interface {
	AddTopicListener(string, CallbackFunc)
	Listen()
	Publish(string, string)
}
