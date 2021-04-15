package common

import (
	"fmt"
	"time"

	"github.com/dxasu/gostar/config"
	"github.com/dxasu/gostar/db"
	log "github.com/dxasu/gostar/util/glog"
)

type REDIS struct {
	db.REDIS
}

func NewRedis(cfgKey string) *REDIS {
	red := db.NewRedis(cfgKey)
	return &REDIS{*red}
}

// DeliverGoodsReq DeliverReq  deliverKey===>>> wait_deliver_mq  {demq}:wait_deliver_mq
func (master *REDIS) DeliverGoodsReq(orderData *OrderInfo) {
	var msg string
	deliverKey := config.GetCfgByKey("deliverKey")
	if deliverKey == "wait_deliver_mq" {
		msg = fmt.Sprintf("DeliverReq(order_id='%s', uid=%d, goods_id=%d, quantity=%d, price=%f, currency='%s', channel=%d)",
			orderData.Order_id, orderData.Uid, orderData.Goods_id, orderData.Deliver_num, orderData.Pay_price, orderData.Currency, orderData.PayType)
	} else {
		msg = fmt.Sprintf("{'uid':%d, 'order_id':'%s', 'goods_id':%d, 'pay_type':%d, 'quantity':%d, 'create_time':%d}",
			orderData.Uid, orderData.Order_id, orderData.Goods_id, orderData.PayType, orderData.Deliver_num, time.Now().Unix())
	}

	retcmd := master.LPush(deliverKey, msg)
	if retcmd.Err() != nil {
		log.Errorf("err DeliverGoodsReq msg:%s err:%+v\n", msg, retcmd.Err())
	} else {
		log.Infof("DeliverGoodsReq retcmd:%+v\n", retcmd)
	}
}
