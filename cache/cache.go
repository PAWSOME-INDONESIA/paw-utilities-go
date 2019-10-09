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
		SetZSet(string, ...redis.Z) error
		GetZSet(string) ([]redis.Z, error)

		HMSetWithExpiration(key string, value map[string]interface{}, ttl time.Duration) error
		HMSet(key string, value map[string]interface{}) error
		HSetWithExpiration(key, field string, value interface{}, ttl time.Duration) error
		HSet(key, field string, value interface{}) error
		HMGet(key string, fields ...string) ([]interface{}, error)
		HGetAll(key string) (map[string]string, error)
		HGet(key, field string, response interface{}) error

		Remove(string) error
		FlushDatabase() error
		FlushAll() error
		Close() error
	}
)
