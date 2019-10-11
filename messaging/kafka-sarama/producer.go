package kafka_sarama

import (
	"time"

	"github.com/Shopify/sarama"
)

func (l *Kafka) Publish(topic, msg string) {
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
