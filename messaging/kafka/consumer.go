package kafka

import (
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/messaging"
)

func (l *Kafka) AddTopicListener(topic string, callback messaging.CallbackFunc) {
	l.mu.Lock()
	defer func() {
		l.mu.Unlock()
	}()
	functions := l.CallbackFunctions[topic]
	functions = append(functions, callback)
	l.CallbackFunctions[topic] = functions
}

func (l *Kafka) Listen() {
	if l.Consumer == nil {
		Log.Info("Cannot Listen to Messaging")
		return
	}

	go func() {
		for err := range l.Consumer.Errors() {
			Log.Infof("Error: %s\n", err.Error())
		}
	}()

	go func() {
		for ntf := range l.Consumer.Notifications() {
			Log.Infof("Rebalanced: %+v\n", ntf)
		}
	}()

	go func() {
		for {
			select {
			case msg, ok := <-l.Consumer.Messages():
				if ok {
					l.Consumer.MarkOffset(msg, "") // mark message as processed
					for _, function := range l.CallbackFunctions[msg.Topic] {
						function(msg.Value)
					}
				}
			}
		}
	}()
}
