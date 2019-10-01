package kafka_sarama

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/logs"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/messaging"
	"sync"
	"time"
)

type (
	Option struct {
		Host             []string
		ConsumerGroup    string
		ProducerRetryMax int
	}
	kafka struct {
		option   Option
		log      logs.Logger
		consumer map[string]sarama.ConsumerGroup
		producer sarama.AsyncProducer
		config   *sarama.Config
		mu       sync.Mutex
	}

	handler struct {
		callbacks []messaging.CallbackFunc
		logger    logs.Logger
	}
)

func (h *handler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *handler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		value := message.Value

		for _, callback := range h.callbacks {
			if err := callback(value); err != nil {
				h.logger.Errorf("failed to proceed message")
			}
		}
		session.MarkMessage(message, "")
	}

	return nil
}

func (k *kafka) ReadWithContext(ctx context.Context, topic string, callbacks []messaging.CallbackFunc) error {
	if len(callbacks) < 1 {
		return errors.New("At least 1 callbacks is required")
	}

	k.mu.Lock()
	if _, ok := k.consumer[topic]; !ok {
		consumer, err := sarama.NewConsumerGroup(k.option.Host, k.option.ConsumerGroup, k.config)
		if err != nil {
			k.mu.Unlock()
			return errors.Wrapf(err, "Failed to Create Consumer Topic %s!", topic)
		}
		k.consumer[topic] = consumer
	}
	k.mu.Unlock()

	consumer := k.consumer[topic]

	handler := handler{
		callbacks: callbacks,
		logger:    k.log,
	}

	if err := consumer.Consume(context.Background(), []string{topic}, &handler); err != nil {
		k.log.Error(err)
		return errors.WithStack(err)
	}
	return nil
}

func (k *kafka) Read(topic string, callback []messaging.CallbackFunc) error {
	return k.ReadWithContext(context.Background(), topic, callback)
}

func (k *kafka) PublishWithContext(ctx context.Context, topic, msg string) error {
	go func(topic, message string) {
		input := k.producer.Input()

		input <- &sarama.ProducerMessage{
			Topic:     topic,
			Key:       nil,
			Value:     sarama.StringEncoder(msg),
			Timestamp: time.Now(),
		}
	}(topic, msg)
	return nil
}

func (k *kafka) Publish(topic, msg string) error {
	return k.PublishWithContext(context.Background(), topic, msg)
}

func (k *kafka) Close() error {
	var err error
	for _, w := range k.consumer {
		if e := w.Close(); e != nil {
			err = e
			k.log.Error(err)
		}
	}
	if err = k.producer.Close(); err != nil {
		return errors.Wrapf(err, "Failed to Close Producer")
	}
	return err
}

func (k *kafka) Ping() error {
	return nil
}

func getOption(option *Option) error {
	if option.ProducerRetryMax == 0 {
		option.ProducerRetryMax = 3
	}
	return nil
}

func New(option Option, log logs.Logger) (messaging.Queue, error) {
	err := getOption(&option)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to Initialize Kafka")
	}

	config := sarama.NewConfig()
	config.Version = sarama.V2_3_0_0
	config.ClientID = option.ConsumerGroup

	// - consumer config
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Return.Errors = true

	// - producer config
	config.Producer.Retry.Max = option.ProducerRetryMax

	producer, err := sarama.NewAsyncProducer(option.Host, config)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to Create Producer! %+v", option)
	}

	return &kafka{
		option:   option,
		config:   config,
		log:      log,
		consumer: make(map[string]sarama.ConsumerGroup),
		producer: producer,
		mu:       sync.Mutex{},
	}, nil
}
