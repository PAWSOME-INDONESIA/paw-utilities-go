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
	NewProducer() error
	AddTopicListener(string, CallbackFunc)
	Listen()
	Publish(string, string)
	Close() error
	CheckSession() (state bool)
}
