package redis_pool

import (
	"github.com/digitalysin/ants"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/cache"
	redis_universal "github.com/tiket/TIX-HOTEL-UTILITIES-GO/cache/redis-universal"
	"sync"
	"time"
)

type (
	PoolCallback func(client cache.Cache)
	Option struct {
		redis_universal.Option
		Pool           int
		NonBlocking    bool
		WorkerDuration time.Duration
	}
	pool struct {
		antsPool  *ants.Pool
		client    cache.Cache
		callbacks []PoolCallback
		option    *Option
		sync.Mutex
	}
)

func New(option *Option) (cache.Cache, error) {
	client, err := redis_universal.New(&option.Option)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	antsPool, err := ants.NewPool(option.Pool, ants.WithNonblocking(option.NonBlocking))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	pool := &pool{
		client:    client,
		antsPool:  antsPool,
		option:    option,
		callbacks: make([]PoolCallback, 0),
	}
	go pool.slave()
	return pool, nil
}

func (p *pool) slave() {
	tick := time.Tick(p.option.WorkerDuration)
	for {
		select {
		case <-tick:
			p.Lock()
			if len(p.callbacks) > 0 {
				free := p.antsPool.Free()
				if free > 0 {
					x := 0
					l := len(p.callbacks)
					if len(p.callbacks) < free {
						x = l
					} else {
						x = free
					}

					cbs := p.callbacks[0:x]
					p.callbacks = p.callbacks[x:]
					for _, cb := range cbs {
						_ = p.antsPool.Submit(func() {
							cb(p.client)
						})
					}
				}
			}
			p.Unlock()
		}
	}
}

func (p *pool) Invoke(callback PoolCallback) {
	p.Lock()
	defer p.Unlock()
	p.callbacks = append(p.callbacks, callback)
}

func (p *pool) Client() cache.Cache {
	return p.client
}

func (p *pool) Ping() error {
	var (
		wg  sync.WaitGroup
		err error
	)

	wg.Add(1)
	p.Invoke(func(client cache.Cache) {
		err = client.Ping()
		wg.Done()
	})
	wg.Wait()
	return err
}

func (p *pool) SetWithExpiration(key string, value interface{}, duration time.Duration) error {
	var (
		wg  sync.WaitGroup
		err error
	)

	wg.Add(1)
	p.Invoke(func(client cache.Cache) {
		err = client.SetWithExpiration(key, value, duration)
		wg.Done()
	})
	wg.Wait()
	return err
}

func (p *pool) Set(key string, value interface{}) error {
	var (
		wg  sync.WaitGroup
		err error
	)

	wg.Add(1)
	p.Invoke(func(client cache.Cache) {
		err = client.Set(key, value)
		wg.Done()
	})
	wg.Wait()
	return err
}

func (p *pool) Get(key string, data interface{}) error {
	var (
		wg  sync.WaitGroup
		err error
	)

	wg.Add(1)
	p.Invoke(func(client cache.Cache) {
		err = client.Get(key, data)
		wg.Done()
	})
	wg.Wait()
	return err
}

func (p *pool) Remove(key string) error {
	var (
		wg  sync.WaitGroup
		err error
	)

	wg.Add(1)
	p.Invoke(func(client cache.Cache) {
		err = client.Remove(key)
		wg.Done()
	})
	wg.Wait()
	return err
}

func (p *pool) FlushDatabase() error {
	var (
		wg  sync.WaitGroup
		err error
	)

	wg.Add(1)
	p.Invoke(func(client cache.Cache) {
		err = client.FlushDatabase()
		wg.Done()
	})
	wg.Wait()
	return err
}

func (p *pool) FlushAll() error {
	var (
		wg  sync.WaitGroup
		err error
	)

	wg.Add(1)
	p.Invoke(func(client cache.Cache) {
		err = client.FlushAll()
		wg.Done()
	})
	wg.Wait()
	return err
}

func (p *pool) Close() error {
	var (
		wg  sync.WaitGroup
		err error
	)

	wg.Add(1)
	p.Invoke(func(client cache.Cache) {
		err = client.Close()
		wg.Done()
	})
	wg.Wait()
	return err
}

func (p *pool) SetZSetWithExpiration(key string, duration time.Duration, data ...redis.Z) error {
	var (
		wg  sync.WaitGroup
		err error
	)

	wg.Add(1)
	p.Invoke(func(client cache.Cache) {
		err = client.SetZSetWithExpiration(key, duration, data...)
		wg.Done()
	})
	wg.Wait()
	return err
}

func (p *pool) SetZSet(key string, data ...redis.Z) error {
	var (
		wg  sync.WaitGroup
		err error
	)

	wg.Add(1)
	p.Invoke(func(client cache.Cache) {
		err = client.SetZSet(key, data...)
		wg.Done()
	})
	wg.Wait()
	return err
}

func (p *pool) GetZSet(key string) ([]redis.Z, error) {
	var (
		wg  sync.WaitGroup
		err error
		data []redis.Z
	)

	wg.Add(1)
	p.Invoke(func(client cache.Cache) {
		data, err = client.GetZSet(key)
		wg.Done()
	})
	wg.Wait()
	return data, err
}

func (p *pool) HMSetWithExpiration(key string, value map[string]interface{}, ttl time.Duration) error {
	var (
		wg  sync.WaitGroup
		err error
	)

	wg.Add(1)
	p.Invoke(func(client cache.Cache) {
		err = client.HMSetWithExpiration(key, value, ttl)
		wg.Done()
	})
	wg.Wait()
	return err
}

func (p *pool) HMSet(key string, value map[string]interface{}) error {
	var (
		wg  sync.WaitGroup
		err error
	)

	wg.Add(1)
	p.Invoke(func(client cache.Cache) {
		err = client.HMSet(key, value)
		wg.Done()
	})
	wg.Wait()
	return err
}

func (p *pool) HSetWithExpiration(key, field string, value interface{}, ttl time.Duration) error {
	var (
		wg  sync.WaitGroup
		err error
	)

	wg.Add(1)
	p.Invoke(func(client cache.Cache) {
		err = client.HSetWithExpiration(key, field, value, ttl)
		wg.Done()
	})
	wg.Wait()
	return err
}

func (p *pool) HSet(key, field string, value interface{}) error {
	var (
		wg  sync.WaitGroup
		err error
	)

	wg.Add(1)
	p.Invoke(func(client cache.Cache) {
		err = client.HSet(key, field, value)
		wg.Done()
	})
	wg.Wait()
	return err
}

func (p *pool) HMGet(key string, fields ...string) ([]interface{}, error) {
	var (
		wg  sync.WaitGroup
		response []interface{}
		err error
	)

	wg.Add(1)
	p.Invoke(func(client cache.Cache) {
		response, err = client.HMGet(key, fields...)
		wg.Done()
	})
	wg.Wait()
	return response, err
}

func (p *pool) HGetAll(key string) (map[string]string, error) {
	var (
		wg  sync.WaitGroup
		response map[string]string
		err error
	)

	wg.Add(1)
	p.Invoke(func(client cache.Cache) {
		response, err = client.HGetAll(key)
		wg.Done()
	})
	wg.Wait()
	return response, err
}

func (p *pool) HGet(key, field string, response interface{}) error {
	var (
		wg  sync.WaitGroup
		err error
	)

	wg.Add(1)
	p.Invoke(func(client cache.Cache) {
		err = client.HGet(key, field, response)
		wg.Done()
	})
	wg.Wait()
	return err
}

func (p *pool) Pipeline() cache.Pipe {
	var (
		wg  sync.WaitGroup
		response cache.Pipe
	)

	wg.Add(1)
	p.Invoke(func(client cache.Cache) {
		response = client.Pipeline()
		wg.Done()
	})
	wg.Wait()
	return response
}
