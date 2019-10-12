package kafka_sarama

import (
	"time"

	"github.com/Shopify/sarama"
)

func (l *Kafka) Publish(topic, msg string) {
	producer, err := sarama.NewAsyncProducerFromClient(l.Client)
	if err != nil {
		return
	}
	defer func() { _ = producer.Close() }()

	producer.Input() <- &sarama.ProducerMessage{
		Topic:     topic,
		Key:       nil,
		Value:     sarama.StringEncoder(msg),
		Timestamp: time.Now(),
	}

	l.Option.Log.Info(topic, msg)
}
