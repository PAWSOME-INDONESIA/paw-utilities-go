package kafka

import (
	"time"

	"github.com/Shopify/sarama"
)

func (l *Kafka) Publish(topic, msg string) {
	Log.Infof("PUBLISH : %s - %s\n", topic, msg)
	go func() {
		input := l.Producer.Input()

		input <- &sarama.ProducerMessage{
			Topic:     topic,
			Key:       nil,
			Value:     sarama.StringEncoder(msg),
			Timestamp: time.Now(),
		}
	}()
}
