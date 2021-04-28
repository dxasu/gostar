package config

import (
	"os"

	log "github.com/dxasu/gostar/util/glog"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

//根据配置初始化
func Init() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("InitConfig err:%+v\n", err)
		}
	}()

	bootcfg := viper.GetString("path")
	viperInit(bootcfg)

	jsonFile := viper.GetString("maincfg.config")
	viperInit(jsonFile)

	initEnv()

	ViperWatch()

	log.Infoln("Init all config success!")
}

func initEnv() {
	cfgenv := viper.GetStringMapString("environment")
	for k, v := range cfgenv {
		os.Setenv(k, v)
	}
	env := os.Getenv("env")
	gin.SetMode(env)
}

// func loadJsonConfig(cfgfile string) string {
// 	filePtr, err := os.Open(cfgfile)
// 	if err != nil {
// 		log.Fatalf("Open file failed err:%+v", err)
// 		return ""
// 	}
// 	defer filePtr.Close()
// 	byteValue, _ := ioutil.ReadAll(filePtr)
// 	config := gjson.Get(string(byteValue), cfgenv).String()
// 	return config
// }
