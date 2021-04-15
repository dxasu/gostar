package db

import (
	"github.com/dxasu/gostar/config"
	log "github.com/dxasu/gostar/util/glog"

	"github.com/go-redis/redis"
)

type REDIS struct {
	redis.Client
}

// NewRedis
func newRedis(addr, passw string) *REDIS {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passw, // no password set
		DB:       0,     // use default DB
	})

	log.Infof("Connection to redis success. addr:%s\n", addr)
	return &REDIS{*client}
}

// InitRedis init db by cfgKey redisCfg
func NewRedis(cfgKey string) *REDIS {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("InitRedis redisCfg:%s err:%+v\n", cfgKey, err)
		}
	}()

	redisCfg := config.GetCfgMapStr(cfgKey)

	addr := redisCfg["host"]
	password := redisCfg["password"]
	return newRedis(addr, password)
}
