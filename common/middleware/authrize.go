package middleware

import (
	"fmt"
	"net/http"

	"github.com/dxasu/gostar/common"
	"github.com/dxasu/gostar/common/Const"
	log "github.com/dxasu/gostar/util/glog"

	"github.com/gin-gonic/gin"
)

// respondWithError
func respondWithError(code int, message string, c *gin.Context) {
	c.JSON(code, gin.H{"error": message})
	log.Infof("respondWithError err:%+v \n", message)
	c.Abort()
}

func Authorized(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json;charset=UTF-8")

	log.Infoln("rawQuery:", c.Request.URL)

	orderId := c.Request.FormValue("order_id")
	if orderId == "" {
		respondWithError(http.StatusBadRequest, "order_id required", c)
		return
	}

	orm, _ := c.Get("*DBORM")
	dbase := orm.(*common.DBORM)

	orderData, err := dbase.QueryOrder(orderId)

	if err != nil {
		log.Warningf("Authorized err:%+v\n", err)
		errMsg := fmt.Sprintf(`{"ret_code":%d, "err_desc":"order not found", "orderId":"%s"}`, Const.CODE_ORDERNOTFOUND, orderId)
		respondWithError(http.StatusBadRequest, errMsg, c)
		return
	}

	c.Set("*OrderInfo", orderData)
}

func Authorized_Check(c *gin.Context) {
	order, _ := c.Get("*OrderInfo")
	orderData := order.(*common.OrderInfo)
	if orderData.Status == Const.STATE_ALLSUCCESS {
		errMsg := `{"ret_code":0, "err_desc": "order commplete"}`
		log.Warningf("order.status:%d, orderid:%s", orderData.Status, orderData.Order_id)
		respondWithError(http.StatusBadRequest, errMsg, c)
	}
}
