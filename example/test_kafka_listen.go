package example

import (
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/logs"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/messaging"
	kafka2 "github.com/tiket/TIX-HOTEL-UTILITIES-GO/messaging/kafka"
	"sync"
)

func main() {
	topic := "test"

	option := kafka2.Option{
		Host:          []string{"localhost:9092"},
		ConsumerGroup: "local",
		Interval:      1,
		RequiredAck:   1,
	}
	log, _ := logs.DefaultLog()

	kfk1, err := kafka2.New(option, log)
	if err != nil {
		panic(err)
	}

	kfk2, err := kafka2.New(option, log)
	if err != nil {
		panic(err)
	}

	log.Info("Starting Listening")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		_ = kfk1.Read(topic, []messaging.CallbackFunc{func(bytes []byte) error {
			log.Info(string(bytes))
			return nil
		}})
	}()
	go func() {
		_ = kfk2.Read(topic, []messaging.CallbackFunc{func(bytes []byte) error {
			log.Info(string(bytes))
			return nil
		}})
	}()
	wg.Wait()
}
