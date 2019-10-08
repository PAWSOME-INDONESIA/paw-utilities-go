package cache

import (
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/util"
	"time"

	"github.com/go-redis/redis"
)

type (
	Cache interface {
		util.Ping
		SetWithExpiration(string, interface{}, time.Duration) error
		Set(string, interface{}) error
		Get(string, interface{}) error

		SetZSetWithExpiration(string, time.Duration, ...redis.Z) error
		SetZSetWith(string, ...redis.Z) error
		GetZSet(string) ([]redis.Z, error)

		Remove(string) error
		FlushDatabase() error
		FlushAll() error
		Close() error
	}
)
