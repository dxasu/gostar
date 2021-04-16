package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dxasu/gostar/common"
	"github.com/dxasu/gostar/config"

	log "github.com/dxasu/gostar/util/glog"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbase *common.DBORM
	redis *common.REDIS
)

func main() {
	defer log.Flush()

	config.Init()
	//go grpclib.GrpcInit()
	dbase, redis = common.CommonInit()

	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	log.Warningf("main signal:%+v", <-chSig)
}
