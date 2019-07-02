package kafka

import (
	"sync"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/pkg/errors"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/logs"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/messaging"
)

type Kafka struct {
	Option            *Option
	Consumer          *cluster.Consumer
	Producer          sarama.AsyncProducer
	CallbackFunctions map[string][]messaging.CallbackFunc
	mu                *sync.Mutex
}

type Option struct {
	Host                 []string
	ConsumerWorker       int
	ConsumerGroup        string
	Strategy             cluster.Strategy
	Heartbeat            int
	ProducerMaxBytes     int
	ProducerRetryMax     int
	ProducerRetryBackoff int
	ListTopics           []string
	Log                  logs.Logger
}

var Log logs.Logger

func New(option *Option) (messaging.MQ, error) {
	if option.Log == nil {
		Log, _ = logs.DefaultLog()
	}

	l := Kafka{
		Option:            option,
		CallbackFunctions: make(map[string][]messaging.CallbackFunc),
		mu:                &sync.Mutex{},
	}

	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Group.PartitionStrategy = l.Option.Strategy
	config.Group.Heartbeat.Interval = time.Duration(l.Option.Heartbeat) * time.Second
	brokers := l.Option.Host
	consumer, err := cluster.NewConsumer(
		brokers,
		l.Option.ConsumerGroup,
		l.Option.ListTopics,
		config,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to Create Consumer! %+v", l.Option)
	}
	l.Consumer = consumer

	configProducer := sarama.NewConfig()
	configProducer.Version = sarama.V0_10_0_0
	configProducer.Producer.Return.Errors = true
	configProducer.Producer.Return.Successes = true
	configProducer.Producer.MaxMessageBytes = l.Option.ProducerMaxBytes
	configProducer.Producer.Retry.Max = l.Option.ProducerRetryMax
	configProducer.Producer.Retry.Backoff = time.Duration(l.Option.ProducerRetryBackoff) * time.Second
	producer, err := sarama.NewAsyncProducer(l.Option.Host, configProducer)

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to Create Producer! %+v", l.Option)
	}
	l.Producer = producer

	return &l, nil
}

func (l Kafka) Close() error {
	if err := l.Consumer.Close(); err != nil {
		return errors.Wrapf(err, "Failed to Close Consumer")
	}

	if err := l.Producer.Close(); err != nil {
		return errors.Wrapf(err, "Failed to Close Producer")
	}
	return nil
}
