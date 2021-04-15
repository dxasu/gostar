package common

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dxasu/gostar/common/Const"
	"github.com/dxasu/gostar/config"
	"github.com/dxasu/gostar/db"
	"github.com/dxasu/gostar/util"
	log "github.com/dxasu/gostar/util/glog"

	"github.com/dxasu/gostar/util/gjson"

	"github.com/gin-gonic/gin"
)

type channelInfo struct {
	ChannelId   int32   `json:"channelId"`
	MethodId    string  `json:"methodId"`
	ChannelName string  `json:"channelName"`
	ImgUrl      string  `json:"imgUrl"`
	GoodsList   []Goods `json:"goodsList"`
}

type channelRet struct {
	Code        int32         `json:"code"`
	ChannelList []channelInfo `json:"channelList"`
}

func Handle_goods_list(c *gin.Context) {
	defer HandleDefer(c)

	channelIdStr := config.GetCfgByKey("channelId")
	channelId, _ := strconv.Atoi(channelIdStr)

	orm, _ := c.Get("*DBORM")
	dbase := orm.(*DBORM)

	goodsLs, err := dbase.GetGoodsList(int32(channelId), "")
	if err != nil {
		log.Warningf("db.GetGoodsList channelId:%+v err:%+v\n", channelId, err)
		c.Status(http.StatusInternalServerError)
		return
	}

	chRetJson, err := json.Marshal(goodsLs)
	if err != nil {
		log.Warningf("json.Marshal err:%+v goodsLs:%+v", err, goodsLs)
		return
	}

	c.String(http.StatusOK, string(chRetJson))
}

func Handle_req_order(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	log.Infoln("rawQuery:", c.Request.URL)

	defer HandleDefer(c)

	userIdStr := c.Query("userid")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)

	if err != nil {
		log.Warningf("orderHandle userId:%s miss", userIdStr)
		return
	}

	goodsIdStr := c.Query("goods_id")
	goodsId, err := strconv.ParseInt(goodsIdStr, 10, 64)

	if err != nil {
		log.Warningf("orderHandle goodsId:%s miss", goodsIdStr)
		c.Status(http.StatusBadRequest)
		return
	}

	userInfo := "" //GetUserInfo(userId)

	if userInfo == "" {
		log.Warningf("invalid userId:%d", userId)
		c.Status(http.StatusBadRequest)
		return
	}

	uidStr := gjson.Get(userInfo, "user_basic.uid").String()
	uid, err := strconv.ParseInt(uidStr, 10, 64)

	if uid == 0 {
		log.Warningf("invalid uid:%d userInfo:%+v", uid, userInfo)
		c.Status(http.StatusBadRequest)
		return
	}

	orm, _ := c.Get("*DBORM")
	dbase := orm.(*DBORM)

	country := gjson.Get(userInfo, "user_basic.country").String()

	app_id := config.GetCfgByKey("app_id")

	goods, err := dbase.GetGoods(goodsId, app_id, country)

	if err != nil {
		log.Warningf("invalid goodsId:%d country:%s app_id:%s err:%+v", goodsId, country, app_id, err)
		c.Status(http.StatusBadRequest)
		return
	}

	channelId := int64(1003)

	orderId := util.NewOrderId(uint64(uid), channelId)

	err = dbase.AddOrder(&PayOrder{
		Order_id:      orderId,
		Uid:           uid,
		To_uid:        uid,
		Goods_id:      goodsId,
		Real_country:  country,
		Real_pay_type: channelId,
		Pay_price:     goods.Price,
		Deliver_num:   goods.Units,
	})

	if err != nil {
		log.Errorf("db.Order kInternalError err:%+v", err)
		log.Warningf(`{"ret_code":%d, "err_desc":"%s", "order_id":"%s"}`, Const.CODE_INTERNALERROR, "internal server error", orderId)
		c.Status(http.StatusBadRequest)
		return
	}

	//c.String(http.StatusOK, `{"ret_code":0, "err_desc":"success", "order_id":"%s"}`, orderId)
	c.JSON(http.StatusOK, gin.H{
		"ret_code": 0,
		"err_desc": "success",
		"order_id": orderId,
	})
}

func Handle_query_order(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	log.Infoln("rawQuery:", c.Request.URL)

	defer HandleDefer(c)

	orderId := c.Query("order_id")

	if orderId == "" {
		log.Warningln("order_id is nil")
		FailureResponse(c)
		return
	}

	orm, exsit := c.Get("*DBORM")
	if !exsit {
		log.Warningf(`Handle_query_order *DBORM is nil`)
		FailureResponse(c)
		return
	}
	dbase := orm.(*DBORM)

	orderData, err := dbase.QueryOrder(orderId)

	if err != nil {
		log.Warningln("order_id not found")
		FailureResponse(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret_code": 0,
		"err_desc": "success",
		"order": gin.H{
			"status":      orderData.Status,
			"status_desc": "orderData.Status",
			"order_id":    orderId,
		},
	})
}

// PayOKAndDeliver
func PayOKAndDeliver(c *gin.Context) {
	defer HandleDefer(c)
	orm, _ := c.Get("*DBORM")
	dbase := orm.(*DBORM)

	order, _ := c.Get("*OrderInfo")
	orderData := order.(*OrderInfo)

	log.Warningf("PayOKAndDeliver pay_flow orderData:%+v\n", orderData)

	err := dbase.UpdateByMap("pay_flow", db.K{
		"status": Const.STATE_PAYSUCCESS,
	}, "order_id", orderData.Order_id)
	if err != nil {
		log.Warningf("PayOKAndDeliver UpdateByMap pay_flow orderId:%s err:%+v\n", orderData.Order_id, err)
	}

	redi, _ := c.Get("*REDIS")
	redi.(*REDIS).DeliverGoodsReq(orderData)
}

// FailureResponse
func FailureResponse(c *gin.Context) {
	// failureURL := config.GetCfgByKey("failure")
	// c.Redirect(http.StatusTemporaryRedirect, failureURL)
	c.Header("Content-Type", "text/html")
	c.HTML(http.StatusOK, "failed.html", gin.H{
		"title": "payment failed",
	})
}

// SuccessResponse
func SuccessResponse(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	c.HTML(http.StatusOK, "success.html", gin.H{
		"title": "payment success",
	})
}

// PendingResponse
func PendingResponse(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	c.HTML(http.StatusOK, "pending.html", gin.H{
		"title": "payment pending",
	})
}
