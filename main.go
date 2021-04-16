package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dxasu/gostar/config"

	log "github.com/dxasu/gostar/util/glog"

	_ "github.com/go-sql-driver/mysql"
)

// var (
// 	dbase *common.DBORM
// 	redis *common.REDIS
// )

func main() {
	defer log.Flush()

	config.Init()
	//go grpclib.GrpcInit()
	//dbase, redis = CommonInit()

	chSig := make(chan os.Signal, 1)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	log.Warningf("main signal:%+v", <-chSig)
}

// func CommonInit() (dbase *common.DBORM, redis *common.REDIS) {
// 	redis = common.NewRedis("redisCfg")
// 	dbase = common.NewDB("mysqlCfg", &gorm.Config{
// 		// PrepareStmt: false,
// 		// DryRun:      false,
// 		// PrepareStmt: false,
// 		// DisableAutomaticPing: true,
// 		// TablePrefix: "t_",   // 表名前缀，`User` 的表名应该是 `t_users`
// 		// SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
// 		// SkipDefaultTransaction: true,
// 	})

// 	go routehandle(func(c *gin.Context) {
// 		c.Set("*DBORM", dbase)
// 		c.Set("*REDIS", redis)
// 	})
// 	return dbase, redis
// }

// // routehandle
// func routehandle(initmiddle func(c *gin.Context)) {
// 	router := gin.Default()
// 	router.Use(initmiddle)
// 	router.LoadHTMLGlob("html/*.html")
// 	router.StaticFS("/static", http.Dir("./static"))

// 	router.GET("/form", middleware.SetLanguage, func(c *gin.Context) {
// 		common.CommitInfo(c, "", []common.PropInfo{
// 			{"Name:", "text", "name", "uname"},
// 			{"Email:", "text", "email", "email"},
// 			{"Document:", "text", "document", "123"},
// 		})
// 	})

// 	subrouter := router.Group("/show")
// 	{
// 		subrouter.POST("/create", common.CreateTransaction)
// 		//subrouter.StaticFS("/static", http.Dir("../../static"))
// 	}

// 	serverHost := config.GetCfgByKey("server")
// 	log.Infof("router.Run server host:%s\n", serverHost)

// 	err := router.Run(serverHost)
// 	//router.RunTLS(serverHost, "server.crt", "server.key")

// 	log.Warningf("ListenAndServe closed err:%+v\n", err)
// }
