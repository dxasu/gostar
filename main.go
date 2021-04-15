package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dxasu/gostar/config"
	"github.com/dxasu/gostar/paylogic"

	log "github.com/dxasu/gostar/util/glog"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	defer log.Flush()

	config.Init()
	//go grpclib.GrpcInit()

	paylogic.Init()

	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	log.Warningf("main signal:%+v", <-chSig)
}
