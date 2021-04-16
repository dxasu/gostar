package common

import (
	"fmt"
	"html/template"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/dxasu/gostar/config"
	"github.com/dxasu/gostar/db"
	log "github.com/dxasu/gostar/util/glog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CommonInit() (dbase *DBORM, redis *REDIS) {
	redis = NewRedis("redisCfg")
	dbase = NewDB("mysqlCfg", &gorm.Config{
		// PrepareStmt: false,
		// DryRun:      false,
		// PrepareStmt: false,
		// DisableAutomaticPing: true,
		// TablePrefix: "t_",   // 表名前缀，`User` 的表名应该是 `t_users`
		// SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		// SkipDefaultTransaction: true,
	})

	go routehandle(func(c *gin.Context) {
		c.Set("*DBORM", dbase)
		c.Set("*REDIS", redis)
	})
	return dbase, redis
}

func setLanguage(c *gin.Context) {
	var lang string
	if lang = c.Param("lang"); lang == "" {
		lang = c.Query("lang")
	}
	if len(lang) < 2 {
		lang = "en"
	}

	lang = strings.ToLower(lang)[0:2]
	langinfo := config.GetMapBYKey("language." + lang)

	c.Set("langinfo", langinfo)
}

// routehandle
func routehandle(initmiddle func(c *gin.Context)) {
	router := gin.Default()
	router.Use(initmiddle)
	router.LoadHTMLGlob("html/*.html")
	router.StaticFS("/static", http.Dir("./static"))

	router.GET("/form", setLanguage, func(c *gin.Context) {
		CommitInfo(c, "", []PropInfo{
			{"Name:", "text", "name", "uname"},
			{"Email:", "text", "email", "email"},
			{"Document:", "text", "document", "123"},
		})
	})

	subrouter := router.Group("/show")
	{
		subrouter.POST("/create", CreateTransaction)
		//subrouter.StaticFS("/static", http.Dir("../../static"))
	}

	serverHost := config.GetCfgByKey("server")
	log.Infof("router.Run server host:%s\n", serverHost)

	err := router.Run(serverHost)
	//router.RunTLS(serverHost, "server.crt", "server.key")

	log.Warningf("ListenAndServe closed err:%+v\n", err)
}

// NotifyMsg
func NotifyMsg(c *gin.Context, code int, msg string) {
	log.Warningf("code:%d msg:%s\n", code, msg)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

// HandleDefer
func HandleDefer(c *gin.Context) {
	if err := recover(); err != nil {
		_, file, line, _ := runtime.Caller(3)
		log.Errorf("->%s:%d] URL:%+v err:%+v\n", file, line, c.Request.URL, err)
		FailureResponse(c)
	}
}

// CreateTransaction
// func CreateTransaction(direct_url string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		defer HandleDefer(c)

// 		order, _ := c.Get("*OrderInfo")
// 		orderData := order.(*OrderInfo)

// 		directURL := fmt.Sprintf(direct_url, orderData.Order_id, orderData.PriceDesc, orderData.GoodsDesc)
// 		log.Infoln("redirect url:", directURL)
// 		c.Redirect(http.StatusTemporaryRedirect, directURL)
// 	}
// }

// CreateTransaction
func CreateTransaction(c *gin.Context) {
	defer HandleDefer(c)

	order, _ := c.Get("*OrderInfo")
	orderData := order.(*OrderInfo)
	if config.GetCfgByKey("needuinfo") != "" {
		orm, _ := c.Get("*DBORM")
		dbase := orm.(*DBORM)
		data, err := dbase.GetByMap("user_extra_data", db.K{
			"uid":        orderData.Uid,
			"channel_id": orderData.PayType,
		})

		if err != nil {
			log.Warningf("getUserInfo orderID:%s err:%+v\n", orderData.Order_id, err)
			FailureResponse(c)
			return
		}

		if data["valid"] != "1" {
			CommitInfo(c, orderData.Order_id, []PropInfo{
				{"Name:", "name", data["uname"]},
				{"Email:", "email", data["email"]},
				{"Document:", "document", data["any01"]},
			})
			return
		}
	}
	//CommitPay(c, orderData)
}

func toint(str string) (ret int) {
	ret, _ = strconv.Atoi(str)
	return
}

// SaveUserInfo
func SaveUserInfo(kv KPROP) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer HandleDefer(c)

		order, _ := c.Get("*OrderInfo")
		orderData := order.(*OrderInfo)

		data := db.K{
			"uid":        orderData.Uid,
			"channel_id": orderData.PayType,
		}

		for k, v := range kv {
			val := c.Request.PostFormValue(v)
			if val == "" {
				log.Warningf("SaveUserInfo err PostFormValue:%s is nil", v)
				FailureResponse(c)
				return
			}
			data[k] = val
		}

		log.Infof("SaveUserInfo data:%+v\n", data)

		orm, _ := c.Get("*DBORM")
		dbase := orm.(*DBORM)

		// err := dbase.Table("user_extra_data").Clauses(clause.OnConflict{
		// 	//DoNothing : true,
		// 	UpdateAll: true,
		// }).Create(data).Error
		err := dbase.SaveByMap("user_extra_data", data)
		if err != nil {
			log.Warningf("SaveUserInfo err:%+v", err)
			FailureResponse(c)
			return
		}

		//CommitPay(c, orderData)
	}
}

type PropInfo []interface{}

// instantiate web html
func instantiate(c *gin.Context, tmpl, action, order_id, format string, data []PropInfo) {
	var content string
	for _, k := range data {
		content += fmt.Sprintf(format, k...)
	}
	langinfo, _ := c.Get("langinfo")
	langs := langinfo.(map[string]interface{})
	modle := make(map[string]interface{})
	modle["content"] = template.HTML(content)
	modle["title"] = langs["infoneed"]
	modle["action"] = action
	modle["infoneed"] = langs["infoneed"]
	modle["order_id"] = order_id
	modle["next"] = langs["next"]
	c.Header("Content-Type", "text/html")
	c.HTML(http.StatusOK, tmpl, modle)
}

func CommitInfo(c *gin.Context, order_id string, data []PropInfo) {
	format := `<div class="collect-item">
	<div class="user-collect-tag">
		<span class="nb_txt">%s</span>:<img src="static/images/star.png" height="8px" width="8px">
	</div>
	<input type="%s" name="%s" class="user-collect-input" value="%s" oninput="onInputChange()">
</div>`
	//action := fmt.Sprintf("saveuserinfo?order_id=%s", orderData.Order_id)
	// datas := []PropInfo{
	// 	{"Name", "name", data["uname"]},
	// 	{"Email", "email", data["email"]},
	// 	{"Document", "document", data["any01"]},
	// }
	instantiate(c, "infoneed.html", "saveuserinfo", order_id, format, data)
}
