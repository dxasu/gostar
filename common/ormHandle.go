package common

import (
	"fmt"
	stdLog "log"
	"os"
	"strings"
	"time"

	"github.com/dxasu/gostar/db"
	log "github.com/dxasu/gostar/util/glog"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBORM struct {
	db.ORM
}

func NewDB(cfgKey string, opts ...gorm.Option) *DBORM {
	ormlog := logger.New(
		stdLog.New(os.Stdout, "\r\n", stdLog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info,
			Colorful:      true,
		})
	opts = append(opts, &gorm.Config{Logger: ormlog})
	db, _ := db.NewORM(cfgKey, opts...)
	return &DBORM{*db}
}

// QueryOrder 查询订单
func (master *DBORM) QueryOrder(orderId string) (*OrderInfo, error) {
	data := new(OrderInfo)
	dbret := master.Table("pay_flow").Select(`pay_flow.order_id,
	pay_flow.uid,
	pay_flow.STATUS as status,
	b.price as pay_price,
	b.deliver_num,
	pay_flow.real_pay_type,
	pay_flow.goods_id,
	pay_flow.third_id,
	pay_flow.real_country,
	b.goods_desc,
	b.price_desc,
	b.currency_type `).Joins("join pay_goods_config AS b ON pay_flow.goods_id = b.goods_id ").Where("pay_flow.order_id = ?", orderId).Scan(data)
	return data, dbret.Error
}

// QueryOrderByThirdID 查询订单
func (master *DBORM) QueryOrderByThirdID(thirdID string) (*OrderInfo, error) {
	data := new(OrderInfo)
	dbret := master.Table("pay_flow").Select(`pay_flow.order_id,
	pay_flow.uid,
	pay_flow.STATUS as status,
	b.price as pay_price,
	b.deliver_num,
	pay_flow.real_pay_type,
	pay_flow.goods_id,
	pay_flow.third_id,
	pay_flow.real_country,
	b.goods_desc,
	b.price_desc,
	b.currency_type `).Joins("join pay_goods_config AS b ON pay_flow.goods_id = b.goods_id ").Where("pay_flow.third_id = ?", thirdID).Scan(data)
	return data, dbret.Error
}

// GetUserPayInfo 获取用户信息
func (master *DBORM) GetUserPayInfo(uid, channelId int64) (data db.K, err error) {
	err = master.Table("user_extra_data").Where(db.K{"uid": uid, "channel_id": channelId}).Scan(&data).Error
	log.Infof("data:%+v\n", data)
	return
}

func (master *DBORM) GetGoodsList(channelId int32, method_id string) ([]Goods, error) {
	return nil, nil
}

func (master *DBORM) AddOrder(orderInfo *PayOrder) error { return nil }

func (master *DBORM) GetGoods(goods_id int64, app_id, country string) (*Goods, error) {
	return nil, nil
}

func (master *DBORM) InsertByMap(table string, data map[string]interface{}) error {
	log.Infof("InsertByMap table:%s data:%+v\n", table, data)
	return master.Table(table).Create(&data).Error
}

func (master *DBORM) UpdateByMap(table string, updates map[string]interface{}, keys ...string) error {
	log.Infof("UpdateByMap table:%s updates:%+v keys:%+v\n", table, updates, keys)
	vals := make([]interface{}, 0)
	query := ""
	for _, key := range keys {
		vals = append(vals, updates[key])
		query += fmt.Sprintf("%s = ? AND ", key)
	}
	query = query[:len(query)-4]
	return master.Table(table).Where(query, vals...).Updates(updates).Error
}

// update when insert failure
func (master *DBORM) SaveByMap(table string, data map[string]interface{}, keys ...string) error {
	log.Infof("SaveByMap table:%s data:%+v\n", table, data)
	err := master.Table(table).Create(&data).Error
	if err == nil || !strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
		return err
	}
	return master.UpdateByMap(table, data, keys...)
}

func (master *DBORM) GetByMap(table string, keys map[string]interface{}) (map[string]interface{}, error) {
	log.Infof("GetByMap table:%s keys:%+v\n", table, keys)
	data := make(map[string]interface{})
	err := master.Table(table).Where(keys).Scan(&data).Error
	return data, err
}

func (master *DBORM) DeleteByMap(table string, keys map[string]interface{}) error {
	log.Infof("DeleteByMap table:%s keys:%+v\n", table, keys)
	vals := make([]interface{}, 0)
	query := ""
	for key, val := range keys {
		vals = append(vals, val)
		query += fmt.Sprintf("%s = ? AND ", key)
	}
	query = query[:len(query)-4]
	return master.Table(table).Where(query, vals...).Delete(&keys).Error
}
