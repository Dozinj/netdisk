package initialization

import (
	"github.com/go-redis/redis"
	"go.uber.org/zap"

	"netdisk/global"
)

func Redis()*redis.Client {
	redisCfg := global.Config.Section("redis")
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Key("host").String() + ":" + redisCfg.Key("port").String(),
		Password: "", // no password
		DB:       0,  // use default DB
	})
	pong, err := client.Ping().Result()
	if err != nil {
		global.SugaredLogger.Panic("redis connect ping failed, err:", zap.Any("err", err))
		return nil
	} else {
		global.SugaredLogger.Info("redis connect ping response:", zap.String("pong", pong))
		return client
	}
}
