package redis_pool

import (
	"github.com/digitalysin/ants"
	"github.com/pkg/errors"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/cache"
	redis_universal "github.com/tiket/TIX-HOTEL-UTILITIES-GO/cache/redis-universal"
	"sync"
	"time"
)

type (
	Option struct {
		redis_universal.Option
		Pool           int
		NonBlocking    bool
		WorkerDuration time.Duration
	}
	pool struct {
		antsPool  *ants.Pool
		client    cache.Cache
		callbacks []cache.PoolCallback
		option    *Option
		sync.Mutex
	}
)

func New(option *Option) (cache.Pool, error) {
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
		callbacks: make([]cache.PoolCallback, 0),
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

func (p *pool) Use(callback cache.PoolCallback) {
	p.Lock()
	defer p.Unlock()
	p.callbacks = append(p.callbacks, callback)
}

func (p *pool) Client() cache.Cache {
	return p.client
}

func (p *pool) Close() error {
	p.antsPool.Release()
	return p.client.Close()
}
