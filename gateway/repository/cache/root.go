package cache

import (
	"encoding/json"
	"github.com/04Akaps/gateway_with_kafka.git/config"
	"github.com/04Akaps/gateway_with_kafka.git/gateway/repository/cache/local"
	"github.com/04Akaps/gateway_with_kafka.git/gateway/repository/cache/redis"
	"log"
	"time"
)

type cache struct {
	cfg config.GatewayCfg

	local local.CacheLocalImpl
	redis redis.RemoteImpl
}

type CacheImpl interface {
}

func NewCache(cfg config.GatewayCfg) (CacheImpl, error) {
	c := cache{cfg: cfg}
	var err error

	ttl := time.Duration(3600) * time.Second

	c.local, err = local.NewLocalCache(ttl)

	if err != nil {
		return nil, err
	}

	c.redis = redis.NewRemote(cfg.Redis.DataSource, cfg.Redis.UserName, cfg.Redis.Password, cfg.Redis.DB)

	go c.runSwitch()

	return c, nil
}

func (c *cache) runSwitch() {
	for {
		e := c.redis.Ping()
		if e != nil {
			// TODO -> redis가 망가진 상태에서, local 캐싱을 사용하라는건데..
			// -> redis 데이터를 local로 동기화 할 필요는 없나?? 필요해보이는데
			c.local.Run()
		} else {
			c.local.Stop()
			_ = c.local.Reset()
		}
		log.Println("cache local running status", "isRunning", c.local.IsRunning())

		<-time.Tick(time.Second * 1)
	}
}

func (c *cache) marshal(v any) ([]byte, error) {
	jsonBytes, err := json.Marshal(v)
	if err == nil {
		return jsonBytes, nil
	}

	return nil, err
}
