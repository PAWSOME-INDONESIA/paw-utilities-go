package example

import (
	"encoding/json"
	"time"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/cache"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/logs"
)

func testRedis() {
	log, _ := logs.DefaultLog()
	option := cache.Option{
		Address:      "localhost:6379",
		DB:           0,
		WriteTimeout: time.Duration(1) * time.Second,
		ReadTimeout:  time.Duration(1) * time.Second,
	}

	redis, err := cache.New(&option)
	if err != nil {
		panic(err)
	}

	key := "key"
	data := make(map[string]bool)
	data["1"] = true

	if err = redis.Set(key, MapBool(data)); err != nil {
		log.Error(err)
	}

	var response MapBool
	if err = redis.Get(key, &response); err != nil {
		log.Error(err)
		return
	}

	log.Info(response["1"])
}

type MapBool map[string]bool

func (m MapBool) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MapBool) UnmarshalBinary(msg []byte) error {
	return json.Unmarshal(msg, &m)
}
