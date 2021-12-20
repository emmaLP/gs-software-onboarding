package caching

import (
	"context"
	"fmt"
	"time"

	"github.com/emmaLP/gs-software-onboarding/internal/database"
	commonModel "github.com/emmaLP/gs-software-onboarding/pkg/common/model"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type Client interface {
	GetAll(ctx context.Context) ([]commonModel.Item, error)
	GetStories(ctx context.Context) ([]commonModel.Item, error)
	GetJobs(ctx context.Context) ([]commonModel.Item, error)
}

type itemCache struct {
	cacheClient *cache.Cache
	dbClient    database.Client
	logger      *zap.Logger
	ringClient  *redis.Ring
	ttl         time.Duration
}

type Options func(c *itemCache)

func WithTTL(ttl time.Duration) Options {
	return func(c *itemCache) {
		c.ttl = ttl
	}
}

func New(ctx context.Context, redisAddr string, db database.Client, logger *zap.Logger, opts ...Options) (*itemCache, error) {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"leader": redisAddr,
		},
	})

	cacheClient := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	_, err := ring.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("Unable to connect with redit. %w", err)
	}

	item := &itemCache{
		dbClient:    db,
		cacheClient: cacheClient,
		ringClient:  ring,
		ttl:         5 * time.Minute,
		logger:      logger,
	}

	for _, opt := range opts {
		opt(item)
	}

	return item, nil
}

func (c *itemCache) GetAll(ctx context.Context) ([]commonModel.Item, error) {
	// var items []commonModel.Item

	key := "items:all"
	return c.cacheItem(key, func(*cache.Item) (interface{}, error) {
		c.logger.Info(fmt.Sprintf("%s caching missed. fetching from source", key))
		return c.dbClient.ListJobs(ctx)
	})
	//err := c.cacheClient.Once(&caching.Item{
	//	Key:   key,
	//	Value: &items,
	//	TTL:   c.ttl,
	//	Do: func(*caching.Item) (interface{}, error) {
	//		c.logger.Info(fmt.Sprintf("%s caching missed. fetching from source", key))
	//		return c.dbClient.ListAll(ctx)
	//	},
	//})
	//if err != nil {
	//	return nil, err
	//}
	//
	//return items, nil
}

func (c *itemCache) GetStories(ctx context.Context) ([]commonModel.Item, error) {
	// var items []commonModel.Item

	key := "items:stories"
	return c.cacheItem(key, func(*cache.Item) (interface{}, error) {
		c.logger.Info(fmt.Sprintf("%s caching missed. fetching from source", key))
		return c.dbClient.ListJobs(ctx)
	})
	//err := c.cacheClient.Once(&caching.Item{
	//	Key:   key,
	//	Value: &items,
	//	TTL:   c.ttl,
	//	Do: func(*caching.Item) (interface{}, error) {
	//		c.logger.Info(fmt.Sprintf("%s caching missed. fetching from source", key))
	//		return c.dbClient.ListStories(ctx)
	//	},
	//})
	//if err != nil {
	//	return nil, err
	//}
	//
	//return items, nil
}

func (c *itemCache) GetJobs(ctx context.Context) ([]commonModel.Item, error) {
	// var items []commonModel.Item

	key := "items:jobs"
	//err := c.cacheClient.Once(&caching.Item{
	//	Key:   key,
	//	Value: &items,
	//	TTL:   c.ttl,
	//	Do: func(*caching.Item) (interface{}, error) {
	//		c.logger.Info(fmt.Sprintf("%s caching missed. fetching from source", key))
	//		return c.dbClient.ListJobs(ctx)
	//	},
	//})
	//if err != nil {
	//	return nil, err
	//}
	//
	//return items, nil

	return c.cacheItem(key, func(*cache.Item) (interface{}, error) {
		c.logger.Info(fmt.Sprintf("%s caching missed. fetching from source", key))
		return c.dbClient.ListJobs(ctx)
	})
}

func (c *itemCache) cacheItem(cacheName string, doFunc func(*cache.Item) (interface{}, error)) ([]commonModel.Item, error) {
	var items []commonModel.Item

	err := c.cacheClient.Once(&cache.Item{
		Key:   cacheName,
		Value: &items,
		TTL:   c.ttl,
		Do:    doFunc,
	})
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (c *itemCache) Close() {
	c.ringClient.Close()
}
