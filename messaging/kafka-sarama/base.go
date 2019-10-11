package kafka_sarama

import (
	"sync"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/pkg/errors"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/logs"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/messaging"
)

const (
	DefaultConsumerWorker       = 10
	DefaultStrategy             = cluster.StrategyRoundRobin
	DefaultHeartbeat            = 3
	DefaultProducerMaxBytes     = 1000000
	DefaultProducerRetryMax     = 3
	DefaultProducerRetryBackoff = 100
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

func getOption(option *Option) {
	if option.Log == nil {
		logger, _ := logs.DefaultLog()
		option.Log = logger
	}

	if option.Strategy == "" {
		option.Strategy = DefaultStrategy
	}

	if option.Heartbeat == 0 {
		option.Heartbeat = DefaultHeartbeat
	}

	if option.ConsumerWorker == 0 {
		option.ConsumerWorker = DefaultConsumerWorker
	}

	if option.ProducerMaxBytes == 0 {
		option.ProducerMaxBytes = DefaultProducerMaxBytes
	}

	if option.ProducerRetryMax == 0 {
		option.ProducerRetryMax = DefaultProducerRetryMax
	}

	if option.ProducerRetryBackoff == 0 {
		option.ProducerRetryBackoff = DefaultProducerRetryBackoff
	}
}

func New(option *Option) (messaging.MessagingQueue, error) {
	getOption(option)

	l := Kafka{
		Option:            option,
		CallbackFunctions: make(map[string][]messaging.CallbackFunc),
		mu:                &sync.Mutex{},
	}

	err := l.NewListener(option)
	if err != nil {
		return nil, err
	}

	err = l.NewProducer()

	if err != nil {
		return nil, err
	}

	return &l, nil
}

func (l *Kafka) CheckSession() (state bool) {
	state = true
	if l.Producer == nil || l.Producer.Input() == nil {
		state = false
	}
	return
}

func (l *Kafka) NewListener(option *Option) error {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Group.PartitionStrategy = l.Option.Strategy
	config.Group.Heartbeat.Interval = time.Duration(l.Option.Heartbeat) * time.Second
	brokers := l.Option.Host
	consumer, err := cluster.NewConsumer(brokers, l.Option.ConsumerGroup, l.Option.ListTopics, config)

	if err != nil {
		return errors.Wrapf(err, "Failed to Create Consumer! %+v", l.Option)
	}
	l.Consumer = consumer

	return err
}

func (l *Kafka) NewProducer() error {
	configProducer := sarama.NewConfig()
	configProducer.Version = sarama.V0_10_0_0
	configProducer.Producer.Return.Errors = true
	configProducer.Producer.Return.Successes = true
	configProducer.Producer.MaxMessageBytes = l.Option.ProducerMaxBytes
	configProducer.Producer.Retry.Max = l.Option.ProducerRetryMax
	configProducer.Producer.Retry.Backoff = time.Duration(l.Option.ProducerRetryBackoff) * time.Millisecond
	client, err := sarama.NewClient(l.Option.Host, configProducer)
	l.Producer, err = sarama.NewAsyncProducerFromClient(client)

	if err != nil {
		return errors.Wrapf(err, "Failed to Create Producer! %+v", l.Option)
	}
	//l.Producer = producer

	//defer func() {
	//	if err := l.Producer.Close(); err != nil {
	//		l.Option.Log.Error(err)
	//	}
	//}()
	//err = producer.Close()
	//if err != nil {
	//	return errors.Wrapf(err, "Failed to Close", l.Option)
	//}

	return nil
}

func (l *Kafka) Close() error {
	if err := l.Consumer.Close(); err != nil {
		return errors.Wrapf(err, "Failed to Close Consumer")
	}
	return nil
}
