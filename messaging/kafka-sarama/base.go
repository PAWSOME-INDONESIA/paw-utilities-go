package kafka_sarama

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/logs"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/messaging"
	"sync"
)

type (
	Option struct {
		Host          []string
		ConsumerGroup string
		SaramaVersion [4]uint
	}
	kafka struct {
		option   Option
		log      logs.Logger
		consumer map[string]sarama.ConsumerGroup
		config   *sarama.Config
		mu       sync.Mutex
	}

	handler struct {
		callback messaging.CallbackFunc
		logger   logs.Logger
	}
)

func (h *handler) Setup(session sarama.ConsumerGroupSession) error {
	panic("implement me")
}

func (h *handler) Cleanup(session sarama.ConsumerGroupSession) error {
	panic("implement me")
}

func (h *handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		value := message.Value

		for _, callback := range h.callback {
			go func(cb messaging.CallbackFunc) {
				if err := callback(value); err != nil {
					h.logger.Errorf("failed to proceed message")
				}
			}(callback)
		}
	}
	panic("implement me")
}

func (k *kafka) ReadWithContext(context.Context, string, []messaging.CallbackFunc) error {
	consumer, err := sarama.NewConsumerGroup(k.option.Host, k.option.ConsumerGroup, k.config)
	if err != nil {
		return errors.WithStack(err)
	}

}

func (k *kafka) Read(topic string, callback []messaging.CallbackFunc) error {
	return k.ReadWithContext(context.Background(), topic, callback)
}

func (k *kafka) PublishWithContext(context.Context, string, string) error {
	panic("implement me")
}

func (k *kafka) Publish(string, string) error {
	panic("implement me")
}

func (k *kafka) Close() error {
	panic("implement me")
}

func getOption(option *Option) error {
	return nil
}

func New(option Option, log logs.Logger) (messaging.Queue, error) {
	err := getOption(&option)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to Initialize Kafka")
	}

	config := sarama.NewConfig()
	config.ClientID = option.ConsumerGroup
	config.Consumer.Return.Errors = true
	config.Version = sarama.KafkaVersion{option.SaramaVersion}

	return &kafka{
		option:   option,
		log:      log,
		consumer: make(map[string]sarama.ConsumerGroup),
		config:   config,
		mu:       sync.Mutex{},
	}, nil
}
