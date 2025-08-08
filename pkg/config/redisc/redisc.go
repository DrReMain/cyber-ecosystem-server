package redisc

import (
	"context"
	"crypto/tls"
	"errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

type RedisConfig struct {
	Host     string `json:",env=REDIS_HOST"`
	DB       int    `json:",default=0,env=REDIS_DB"`
	Username string `json:",optional,env=REDIS_USERNAME"`
	Password string `json:",optional,env=REDIS_PASSWORD"`
	TLS      bool   `json:",optional,env=REDIS_TLS"`
	Master   string `json:",optional,env=REDIS_MASTER"`
}

func (c RedisConfig) Validate() error {
	if len(c.Host) == 0 {
		return errors.New("redis host is required")
	}
	return nil
}

func (c RedisConfig) NewUniversalClient() (redis.UniversalClient, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	opt := &redis.UniversalOptions{
		Addrs:    strings.Split(c.Host, ","),
		DB:       c.DB,
		Username: c.Username,
		Password: c.Password,
	}

	if c.Master != "" {
		opt.MasterName = c.Master
	}

	if c.TLS {
		opt.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	client := redis.NewUniversalClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return client, nil
}

func (c RedisConfig) MustNewUniversalClient() redis.UniversalClient {
	client, err := c.NewUniversalClient()
	if err != nil {
		logx.Must(err)
	}
	return client
}
