package db

import (
	log "github.com/dxasu/gostar/util/glog"

	"github.com/go-redis/redis"
)

// GetToken by key
func (redisHandle *REDIS) GetToken(key string) string {
	ret := redisHandle.Get(key)
	val, err := ret.Result()
	if err == redis.Nil {
		return ""
	} else if err != nil {
		log.Errorf("redis err as GetToken key:%s", key)
		return ""
	} else {
		return val
	}
}

// func handle_task(req *redis.StringCmd) {

// 	config.CompleteOrder()
// }

// func Dowork() {
// 	for {
// 		req := redisHandle.BRPopLPush("wait_deliver_mq", "do_deliver_mq", 30)
// 		if req != nil {
// 			log.Infof("recv new task:%+v", *req)
// 			handle_task(req)
// 			redisHandle.LRem("do_deliver_mq", 0, req)
// 		}
// 	}
// }
