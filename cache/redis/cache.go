package redis

import (
	"encoding"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/cache"
	"time"
)

type (
	Option struct {
		Address      string
		Password     string
		DB           int
		PoolSize     int
		MinIdleConns int
		DialTimeout  time.Duration
		PoolTimeout  time.Duration
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		MaxConnAge   time.Duration
	}

	redisClient struct {
		r *redis.Client
	}
)

func New(option *Option) (cache.Cache, error) {
	var client *redis.Client

	client = redis.NewClient(&redis.Options{
		DB:           option.DB,
		Addr:         option.Address,
		Password:     option.Password,
		PoolSize:     option.PoolSize,
		PoolTimeout:  option.PoolTimeout,
		ReadTimeout:  option.ReadTimeout,
		WriteTimeout: option.WriteTimeout,
		DialTimeout:  option.DialTimeout,
		MinIdleConns: option.MinIdleConns,
		MaxConnAge:   option.MaxConnAge,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, errors.Wrap(err, "Failed to connect to redis!")
	}

	return &redisClient{r: client}, nil
}

func (c *redisClient) Ping() error {
	if _, err := c.r.Ping().Result(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *redisClient) SetWithExpiration(key string, value interface{}, duration time.Duration) error {
	if err := check(c); err != nil {
		return err
	}

	if _, err := c.r.Set(key, value, duration).Result(); err != nil {
		return errors.Wrapf(err, "failed to set cache with key %s!", key)
	}

	return nil
}

func (c *redisClient) Set(key string, value interface{}) error {
	if err := check(c); err != nil {
		return err
	}

	return c.SetWithExpiration(key, value, 0)
}

func (c *redisClient) Get(key string, data interface{}) error {
	if _, ok := data.(encoding.BinaryUnmarshaler); !ok {
		return errors.New(fmt.Sprintf("failed to get cache with key %s!: redis: can't unmarshal (implement encoding.BinaryUnmarshaler)", key))
	}

	if err := check(c); err != nil {
		return err
	}

	val, err := c.r.Get(key).Result()

	if err == redis.Nil {
		return errors.Wrapf(err, "key %s does not exits", key)
	}

	if err != nil {
		return errors.Wrapf(err, "failed to get key %s!", key)
	}

	if err := data.(encoding.BinaryUnmarshaler).UnmarshalBinary([]byte(val)); err != nil {
		return err
	}

	return nil
}

func (c *redisClient) Remove(key string) error {
	if err := check(c); err != nil {
		return err
	}

	if _, err := c.r.Del(key).Result(); err != nil {
		return errors.Wrapf(err, "failed to remove key %s!", key)
	}

	return nil
}

func (c *redisClient) FlushDatabase() error {
	if err := check(c); err != nil {
		return err
	}

	if _, err := c.r.FlushDB().Result(); err != nil {
		return errors.Wrap(err, "failed to flush db!")
	}

	return nil
}

func (c *redisClient) FlushAll() error {
	if err := check(c); err != nil {
		return err
	}

	if _, err := c.r.FlushAll().Result(); err != nil {
		return errors.Wrap(err, "failed to flush db!")
	}

	return nil
}

func (c *redisClient) Close() error {
	if err := c.r.Close(); err != nil {
		return errors.Wrap(err, "failed to close redis client")
	}

	return nil
}

func check(c *redisClient) error {
	if c.r == nil {
		return errors.New("redis client is not connected")
	}

	return nil
}

func (c *redisClient) SetZSetWithExpiration(key string, duration time.Duration, data ...redis.Z) error {
	if err := c.SetZSetWith(key, data...); err != nil {
		return err
	}

	if _, err := c.r.Expire(key, duration).Result(); err != nil {
		return errors.Wrapf(err, "failed to zadd cache with key %s!", key)
	}
	return nil
}

func (c *redisClient) SetZSetWith(key string, data ...redis.Z) error {
	if err := check(c); err != nil {
		return err
	}

	c.r.Del(key)
	if _, err := c.r.ZAdd(key, data...).Result(); err != nil {
		return errors.Wrapf(err, "failed to zadd cache with key %s!", key)
	}
	return nil
}

func (c *redisClient) GetZSet(key string) ([]redis.Z, error) {
	data, err := c.r.ZRangeWithScores(key, 0, -1).Result()
	if err != nil {
		return nil, errors.Wrap(err, "failed to run zrange command")
	}

	if len(data) <= 0 {
		return nil, errors.New(fmt.Sprintf("key %s does not exits", key))
	}

	return data, nil
}
