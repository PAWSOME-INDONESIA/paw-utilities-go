package cache

import (
	"time"

	"github.com/pkg/errors"

	"github.com/go-redis/redis"
)

type (
	Option struct {
		Address      string
		Password     string
		DB           int
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}

	Cache interface {
		Set(string, interface{}, time.Duration) error
		Get(string) (string, error)
		Remove(string) error
		FlushDatabase() error
	}

	cache struct {
		r *redis.Client
	}
)

func New(option *Option) (Cache, error) {
	var client *redis.Client

	client = redis.NewClient(&redis.Options{
		Addr:         option.Address,
		Password:     option.Password,
		DB:           option.DB,
		ReadTimeout:  option.ReadTimeout,
		WriteTimeout: option.WriteTimeout,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, errors.Wrap(err, "Failed to connect to redis!")
	}

	return &cache{r: client}, nil
}

func (c *cache) Set(key string, value interface{}, duration time.Duration) error {
	if err := check(c); err != nil {
		return err
	}

	if _, err := c.r.Set(key, value, duration).Result(); err != nil {
		return errors.Wrapf(err, "failed to set cache with key %s!", key)
	}

	return nil
}

func (c *cache) Get(key string) (string, error) {
	if err := check(c); err != nil {
		return "", err
	}

	val, err := c.r.Get(key).Result()

	if err == redis.Nil {
		return "", errors.Wrapf(err, "key %s does not exits", key)
	}

	if err != nil {
		return "", errors.Wrapf(err, "failed to get key %s!", key)
	}

	return val, nil
}

func (c *cache) Remove(key string) error {
	if err := check(c); err != nil {
		return err
	}

	if _, err := c.r.Del(key).Result(); err != nil {
		return errors.Wrapf(err, "failed to remove key %s!", key)
	}

	return nil
}

func (c *cache) FlushDatabase() error {
	if err := check(c); err != nil {
		return err
	}

	if _, err := c.r.FlushDB().Result(); err != nil {
		return errors.Wrap(err, "failed to flush db!")
	}

	return nil
}

func check(c *cache) error {
	if c.r == nil {
		return errors.New("redis client is not connected")
	}

	return nil
}
