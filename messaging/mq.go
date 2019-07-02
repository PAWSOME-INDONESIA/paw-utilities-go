package messaging

type CallbackFunc func([]byte) error

type MQ interface {
	Publish(string, string)
	AddTopicListener(string, CallbackFunc)
	Listen()
	Close() error
}
