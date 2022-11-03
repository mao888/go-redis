package gredis

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type GRedis struct {
	client redis.UniversalClient
	config configuration
}
type Message struct {
	Channel      string
	Pattern      string
	Payload      string
	PayloadSlice []string
}

var (
	ErrNotFound  = redis.Nil
	ErrConnect   = errors.New("connect error")
	ErrRedisType = errors.New("redis type error")

	lock sync.Mutex

	redisMap = make(map[string]*GRedis)
)

func MustInit(name string, opts ...Option) {
	conf := _defaultConfiguration
	for _, opt := range opts {
		opt.apply(&conf)
	}
	if len(conf.RedisType) == 0 {
		conf.RedisType = Client
	}
	if _client, err := selfInit(conf); err != nil {
		panic(err)
	} else {
		lock.Lock()
		defer lock.Unlock()
		redisMap[name] = &GRedis{client: _client, config: conf}
	}
}

func Redis(name string) *GRedis {
	return redisMap[name]
}

func RedisClient(name string) redis.UniversalClient {
	gRedis := redisMap[name]
	if gRedis == nil {
		return nil
	}
	return gRedis.client
}

//selfInit Redis初始化连接
func selfInit(conf configuration) (redis.UniversalClient, error) {
	var _client redis.UniversalClient
	switch conf.RedisType {
	case Client:
		addr := ""
		if len(conf.Hosts) > 0 {
			addr = conf.Hosts[0]
		}
		_client = redis.NewClient(&redis.Options{
			Addr:         addr,
			Password:     conf.PassWord,
			PoolSize:     conf.PoolSize,
			MinIdleConns: conf.MinIdleCons,
		})
	case Sentinel:
		_client = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    conf.MasterName,
			SentinelAddrs: conf.Hosts,
			Password:      conf.PassWord,
			PoolSize:      conf.PoolSize,
			MinIdleConns:  conf.MinIdleCons,
		})
	case Cluster:
		_client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        conf.Hosts,
			Password:     conf.PassWord,
			PoolSize:     conf.PoolSize,
			MinIdleConns: conf.MinIdleCons,
		})
	default:
		return nil, ErrRedisType
	}
	if _client == nil {
		return nil, ErrConnect
	}
	if _, err := _client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	return _client, nil
}

func (g *GRedis) Ping() error {
	return g.client.Ping(context.Background()).Err()
}

func (g *GRedis) Get(ctx context.Context, key string) (string, error) {
	return g.client.Get(ctx, key).Result()
}

func (g *GRedis) Set(ctx context.Context, key string, val string, expired time.Duration) error {
	return g.client.Set(ctx, key, val, expired).Err()
}

func (g *GRedis) HGet(ctx context.Context, key, field string) (string, error) {
	return g.client.HGet(ctx, key, field).Result()
}

func (g *GRedis) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return g.client.HGetAll(ctx, key).Result()
}

func (g *GRedis) HSet(ctx context.Context, key, field, val string) error {
	return g.client.HSet(ctx, key, field, val).Err()
}

func (g *GRedis) HDel(ctx context.Context, key string, field ...string) error {
	return g.client.HDel(ctx, key, field...).Err()
}

func (g *GRedis) HMSet(ctx context.Context, key string, data map[string]interface{}) error {
	return g.client.HMSet(ctx, key, data).Err()
}

func (g *GRedis) HIncrBy(ctx context.Context, key, filed string, incr int64) (int64, error) {
	return g.client.HIncrBy(ctx, key, filed, incr).Result()
}

func (g *GRedis) HIncrByFloat(ctx context.Context, key, filed string, incr float64) (float64, error) {
	return g.client.HIncrByFloat(ctx, key, filed, incr).Result()
}

func (g *GRedis) SMembers(ctx context.Context, key string) ([]string, error) {
	return g.client.SMembers(ctx, key).Result()
}

func (g *GRedis) Keys(ctx context.Context, key string) (keys []string) {
	keys, _ = g.KeysE(ctx, key)
	return
}

func (g *GRedis) KeysE(ctx context.Context, key string) (keys []string, err error) {
	var cursor uint64
	var tempKeys []string
	for {
		tempKeys, cursor, err = g.client.Scan(ctx, cursor, key, 10).Result()
		if err != nil {
			break
		}
		keys = append(keys, tempKeys...)
		if cursor == 0 {
			break
		}
	}
	return
}

func (g *GRedis) SetKeys(ctx context.Context, key string) (values []string) {
	var index uint64
	var tempValues []string
	var err error
	for {
		tempValues, index, err = g.client.SScan(ctx, key, index, "", -1).Result()
		if err != nil {
			break
		}
		values = append(values, tempValues...)
		if index == 0 {
			break
		}
	}
	return
}

func (g *GRedis) Expire(ctx context.Context, key string, duration time.Duration) error {
	return g.client.Expire(ctx, key, duration).Err()
}

func (g *GRedis) Del(ctx context.Context, key ...string) error {
	return g.client.Del(ctx, key...).Err()
}

func (g *GRedis) SAdd(ctx context.Context, key, value string) error {
	return g.client.SAdd(ctx, key, value).Err()
}

func (g *GRedis) SScan(ctx context.Context,
	key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return g.client.SScan(ctx, key, cursor, match, count).Result()
}

func (g *GRedis) Subscribe(ctx context.Context, chanNames ...string) (bool, <-chan *redis.Message) {
	pubSub := g.client.Subscribe(ctx, chanNames...)
	if _, err := pubSub.Receive(ctx); err != nil {
		return false, nil
	}
	////redis 默认是100，此处要和Redis默认值保持一致
	//c := make(chan *Message, 100)
	//for m := range pubSub.Channel() {
	//	c <- &Message{
	//		Channel:      m.Channel,
	//		Pattern:      m.Pattern,
	//		Payload:      m.Payload,
	//		PayloadSlice: m.PayloadSlice,
	//	}
	//}
	return true, pubSub.Channel()
}

func (g *GRedis) Publish(ctx context.Context, key, message string) (int64, error) {
	return g.client.Publish(ctx, key, message).Result()
}

func (g *GRedis) Lock(ctx context.Context, key string, expiration time.Duration) bool {
	lock, err := g.client.SetNX(ctx, key, "", expiration).Result()
	if err != nil {
		return false
	}
	return lock
}

func (g *GRedis) UnLock(ctx context.Context, key string) {
	_ = g.Del(ctx, key)
}
