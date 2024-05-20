package redis

import (
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/redis/go-redis/v9"
	"timeLineGin/pkg/config"
)

var instance *redis.Client

func GetInstance() *redis.Client {
	return instance
}

func Initialize(conf *config.Redis) {
	ins := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, gconv.Int(conf.Port)),
		Password: conf.Pass,
		DB:       0,
		PoolSize: 10,
	})
	instance = ins

	//ins.Set(context.Background(), "test", "test", 0)
}
