package common

import (
	"database/sql"
)

type KPROP map[string]string

type ChannelInfo struct {
	Channel_id    int32
	Method_id     string
	Region        string
	Icon          string
	Internal_name string
	Discount_info string
}

type PayData struct {
	Order_id      sql.NullString `db:"order_id"`
	Uid           sql.NullString `db:"uid"`
	Goods_id      sql.NullString `db:"goods_id"`
	Pay_price     sql.NullString `db:"pay_price"`
	Deliver_num   sql.NullString `db:"deliver_num"`
	Real_pay_type sql.NullString `db:"real_pay_type"`
	Country       sql.NullString `db:"country"`
	Price         float64        `db:"price"`
	ItemName      sql.NullString `db:"goods_desc"`
}

// 获取商品信息
type Goods struct {
	Goods_id    int32
	Goods_desc  string
	Price       float64
	Price_desc  string
	ChProductId string `db:"channel_product_id"`
	Units       int32
}

// PayOrder
type PayOrder struct {
	Order_id      string
	Uid           int64
	To_uid        int64
	Goods_id      int64
	Real_country  string
	Real_pay_type int64
	Pay_price     float64
	Deliver_num   int32
}

type OrderInfo struct {
	Order_id    string
	Uid         int64
	Status      int32
	Pay_price   float64
	Deliver_num int64
	PayType     int32 `gorm:"column:real_pay_type"`
	Goods_id    int64
	Third_id    sql.NullString
	Country     string `gorm:"column:real_country"`
	Currency    string `gorm:"column:currency_type"`
	GoodsDesc   string `gorm:"column:goods_desc"`
	PriceDesc   string `gorm:"column:price_desc"`
}
