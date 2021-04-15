package common

import (
	"fmt"

	"github.com/dxasu/gostar/db"
	log "github.com/dxasu/gostar/util/glog"
)

var (
	regionMap map[string]string
)

type DBASE struct {
	db.SQLX
}

func getRegion(country string) string {
	if toRegion, ok := regionMap[country]; ok {
		return fmt.Sprintf("'%s','%s','AA'", country, toRegion)
	}

	return fmt.Sprintf("'%s','AA'", country)
}

func (this *DBASE) GetGoods(goods_id int64, app_id, country string) (*Goods, error) {
	sql := fmt.Sprintf(`SELECT
			price,
			units 
		FROM
			goods_on_shelf 
		WHERE
			goods_id = ? 
			AND app_id = ?
			AND region IN ( %s ) 
			LIMIT 1`, getRegion(country))
	data := Goods{}
	err := this.Get(&data, sql, goods_id, app_id)

	return &data, err
}

// get PayData by orderId
func (this *DBASE) QueryPay(orderId string) (*PayData, error) {
	sql := `SELECT
		pay_flow.order_id,
		pay_flow.uid,
		pay_flow.goods_id,
		pay_flow.pay_price,
		pay_flow.deliver_num,
		pay_flow.real_pay_type,
		pay_goods_config.country,
		pay_goods_config.price,
		pay_goods_config.goods_desc
	FROM
		pay_flow
		LEFT JOIN pay_goods_config ON pay_flow.goods_id = pay_goods_config.goods_id 
	where
		pay_flow.order_id = '?'
		LIMIT 1`

	var data *PayData
	err := this.Get(data, sql, orderId)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (this *DBASE) GetGoodsList(channelId int32, method_id string) ([]Goods, error) {
	sql := fmt.Sprintf(`SELECT
			goods_id,
			goods_desc,
			price,
			price_desc,
			"productId" as channel_product_id,
			units 
		FROM
			goods_on_shelf 
		WHERE
			goods_id IN ( SELECT goods_id FROM channel_goods AS b WHERE b.channel_id = %d AND b.method_id = '%s' ) 
		ORDER BY
			pos_num 
			AND show_mask > 0`, channelId, method_id)
	data := make([]Goods, 0, 8)
	err := this.Select(&data, sql)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (this *DBASE) GetChannelName(channelId string) ([]Goods, error) {
	sql := `SELECT
		goods_on_shelf.goods_id,
		price,
		units,
		goods_desc,
		pay_method_basic_info.internal_name 
	FROM
		goods_on_shelf
		inner JOIN goods_basic_info ON goods_on_shelf.goods_id = goods_basic_info.goods_id
		inner JOIN pay_method_basic_info ON pay_method_basic_info.channel_id = goods_basic_info.channel_id 
	WHERE
		goods_basic_info.channel_id = ?`
	data := make([]Goods, 0, 8)
	err := this.Select(&data, sql, channelId)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetChannelLs grade roleId
func (this *DBASE) GetChannelLs(grade, roleId int64, country, app_id string) ([]ChannelInfo, error) {
	sql := fmt.Sprintf(`SELECT
			a.channel_id, a.method_id, region, icon, internal_name, discount_info
		FROM
			pay_method_on_shelf a
			INNER JOIN pay_method_basic_info b ON a.channel_id = b.channel_id 
			AND a.method_id = b.method_id AND b.enabled = 1 AND (b.support_role = 3 or b.support_role = %d)
		WHERE
			a.min_level <= %d 
			AND a.app_id = '%s'
			AND a.region IN ( %s ) 
		ORDER BY
		b.sort_weight`, roleId, grade, app_id, getRegion(country))
	data := make([]ChannelInfo, 0, 4)
	err := this.Select(&data, sql)
	if err != nil {
		return nil, err
	}
	log.Warningf("%+v\n", data)
	return data, nil
}

// AddOrder 下订单
func (master *DBASE) AddOrder(orderInfo *PayOrder) error {
	sql := `insert into pay_flow (order_id, uid, to_uid, goods_id, status, real_country, real_pay_type, pay_price, deliver_num) 
		values (?, ?, ?, ?, 1, ?, ?, ?, ?)`
	_, err := master.Exec(sql, orderInfo.Order_id, orderInfo.Uid, orderInfo.To_uid, orderInfo.Goods_id, orderInfo.Real_country,
		orderInfo.Real_pay_type, orderInfo.Pay_price, orderInfo.Deliver_num)
	return err
}

// GetUserPayInfo 获取用户信息
func (master *DBASE) GetUserPayInfo(uid, channelId int64) (map[string]string, error) {
	return master.GetByMap("user_extra_data", db.K{"uid": uid, "channel_id": channelId})
}

// yoho sql
// SELECT
// 		a.order_id,
// 		a.uid,
// 		a.STATUS as status,
// 		b.pay_money AS pay_price,
// 		a.deliver_num,
// 		a.real_pay_type,
// 		a.goods_id,
// 		a.third_id,
// 		a.real_country,
// 		b.des AS goods_desc,
// 		b.price AS price_desc,
// 		b.currency AS currency_type
// 	FROM
// 		pay_flow AS a
// 		JOIN pay_goods_config AS b ON a.goods_id = b.goods_id
// 	WHERE
// 		a.order_id = '%s'
// 		LIMIT 1

// QueryOrder 查询订单
func (master *DBASE) QueryOrder(orderId string) (*OrderInfo, error) {
	sqlStr := `SELECT
		a.order_id,
		a.uid,
		a.STATUS as status,
		b.price as pay_price,
		b.deliver_num,
		a.real_pay_type,
		a.goods_id,
		a.third_id,
		a.real_country,
		b.goods_desc,
		b.price_desc,
		b.currency_type 
	FROM
		pay_flow AS a
		JOIN pay_goods_config AS b ON a.goods_id = b.goods_id 
	WHERE
		a.order_id = '%s' 
		LIMIT 1`
	sqlStr = fmt.Sprintf(sqlStr, orderId)
	data := new(OrderInfo)
	err := master.Get(data, sqlStr)
	return data, err
}

// QueryOrderByThirdID 查询订单
func (master *DBASE) QueryOrderByThirdID(thirdID string) (*OrderInfo, error) {
	sqlStr := `SELECT
		a.order_id,
		a.uid,
		a.STATUS as status,
		b.price as pay_price,
		b.deliver_num,
		a.real_pay_type,
		a.goods_id,
		a.third_id,
		a.real_country,
		b.goods_desc,
		b.price_desc,
		b.currency_type 
	FROM
		pay_flow AS a
		JOIN pay_goods_config AS b ON a.goods_id = b.goods_id 
	WHERE
		a.third_id = '%s' 
		LIMIT 1`
	sqlStr = fmt.Sprintf(sqlStr, thirdID)
	data := new(OrderInfo)
	err := master.Get(data, sqlStr)
	return data, err
}
