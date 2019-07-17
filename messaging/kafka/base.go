package kafka

import (
	"context"
	"github.com/pkg/errors"
	kfk "github.com/segmentio/kafka-go"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/logs"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/messaging"
	"time"
)

type (
	kafka struct {
		option Option
		log    logs.Logger
	}
)

type Option struct {
	Host          []string
	ConsumerGroup string
	Interval      int
}

func getOption(option *Option) error {
	if len(option.Host) == 0 {
		return errors.New("Host is required!")
	}
	if option.ConsumerGroup == "" {
		return errors.New("ConsumerGroup is required!")
	}
	if option.Interval == 0 {
		option.Interval = 1
	}
	return nil
}

func New(option Option, log logs.Logger) (messaging.Queue, error) {
	err := getOption(&option)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to Initialize Kafka")
	}
	return &kafka{
		option: option,
		log:    log,
	}, nil
}

func (k *kafka) ReadWithContext(ctx context.Context, topic string, callbacks []messaging.CallbackFunc) error {
	if len(callbacks) < 1 {
		return errors.New("At least 1 callbacks is required")
	}

	reader := kfk.NewReader(kfk.ReaderConfig{
		Brokers: k.option.Host,
		GroupID: k.option.ConsumerGroup,
		Topic:   topic,
		MaxWait: time.Duration(k.option.Interval) * time.Millisecond,
	})
	defer func() {
		if err := reader.Close(); err != nil {
			k.log.Error(err)
		}
	}()

	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			k.log.Error(err)
			continue
		}

		for _, c := range callbacks {
			if err = c(m.Value); err != nil {
				k.log.Error(err)
			}
		}
	}
}

func (k *kafka) Read(topic string, callbacks []messaging.CallbackFunc) error {
	return k.ReadWithContext(context.Background(), topic, callbacks)
}

func (k *kafka) PublishWithContext(ctx context.Context, topic, message string) error {
	w := kfk.NewWriter(kfk.WriterConfig{
		Brokers:      k.option.Host,
		Topic:        topic,
		Balancer:     &kfk.Hash{},
		BatchTimeout: time.Duration(k.option.Interval) * time.Millisecond,
	})
	defer func() {
		if err := w.Close(); err != nil {
			k.log.Error(err)
		}
	}()

	if err := w.WriteMessages(context.Background(), kfk.Message{Value: []byte(message)}); err != nil {
		k.log.Error(err)
		return errors.Wrapf(err, "failed to publish message on topic %s", topic)
	}
	return nil
}

func (k *kafka) Publish(topic, message string) error {
	return k.PublishWithContext(context.Background(), topic, message)
}
