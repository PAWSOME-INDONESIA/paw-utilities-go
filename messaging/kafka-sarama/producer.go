package kafka_sarama

import (
	"time"

	"github.com/Shopify/sarama"
)

func (l *Kafka) Publish(topic, msg string) {
	if err := l.NewProducer(); err != nil {
		l.Option.Log.Error(err)
	}

	defer func() {
		if l.Producer != nil {
			if err := l.Producer.Close(); err != nil {
				l.Option.Log.Error(err)
			}
		}
	}()

	if l.Producer == nil || l.Producer.Input() == nil {
		l.Option.Log.Error("error nil")
		return
	}

	l.Producer.Input() <- &sarama.ProducerMessage{
		Topic:     topic,
		Key:       nil,
		Value:     sarama.StringEncoder(msg),
		Timestamp: time.Now(),
	}
	l.Option.Log.Info("haha ", topic, msg)

}
