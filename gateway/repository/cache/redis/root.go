package redis

import (
	"github.com/go-redis/redis/v7"
	"time"
)

type remote struct {
	client *redis.Client
}

type RemoteImpl interface {
	Ping() error
	Set(key string, value []byte, ttl time.Duration) error
	SetNX(key string, value []byte, ttl time.Duration) error
	Get(key string) ([]byte, error)
	Delete(key string) error
}

func NewRemote(dataSource, password, userName string, db int) RemoteImpl {
	client := redis.NewClient(&redis.Options{
		Addr:     dataSource,
		Password: password,
		DB:       db,
		Username: userName,
	})

	return &remote{
		client: client,
	}
}

func (r *remote) Ping() error {
	return r.client.Ping().Err()
}

func (r *remote) Set(key string, value []byte, ttl time.Duration) error {
	return r.client.Set(key, value, ttl).Err()
}

func (r *remote) SetNX(key string, value []byte, ttl time.Duration) error {
	return r.client.SetNX(key, value, ttl).Err()
}

func (r *remote) Get(key string) ([]byte, error) {
	return r.client.Get(key).Bytes()
}

func (r *remote) Increment(key string) (int64, error) {
	return r.client.Incr(key).Result()
}

func (r *remote) Delete(key string) error {
	return r.client.Del(key).Err()
}
