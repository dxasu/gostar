package config

import (
	"fmt"

	log "github.com/dxasu/gostar/util/glog"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Address := viper.GetString("Address")
// key取Address或者address都能取到值，viper转成小写处理
// viper.SetDefault("LayoutDir", "layouts")
// viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})
// viperInit config.json
func viperInit(cfgfile string) {
	viper.SetConfigFile(cfgfile)
	err := viper.ReadInConfig()

	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		panic(fmt.Errorf("fatal miss config file: %+v", err))
	} else if err != nil {
		panic(fmt.Errorf("fatal error config file: %+v", err))
	}

	//log.Infof("%+v", viper.AllSettings())
}

// GetCfgByKey
func GetCfgByKey(key string) string {
	return viper.GetString(key)
}

// GetByKey
func GetByKey(key string) interface{} {
	return viper.Get(key)
}

// GetMapBYKey 获取配置信息
func GetMapBYKey(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}

// GetCfgMapStr 获取配置信息
func GetCfgMapStr(key string) map[string]string {
	return viper.GetStringMapString(key)
}

// ViperWatch
func ViperWatch() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infoln("config vary:", e.Name)
		log.Infof("%+v", viper.AllKeys())
	})
}

// AddWatch
func AddWatch(confFunc func()) {
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infoln("config vary:", e.Name, e.Op.String())
		log.Infof("%+v", viper.AllKeys())
		confFunc()
	})
}

// GetEnvInfo 获取环境变量
func GetEnvInfo(env string) string {
	viper.AutomaticEnv()
	return viper.GetString(env)
}

func viperEtcd() {
	viper.AddRemoteProvider("etcd",
		"http://127.0.0.1:2379",
		"conf.toml")
	viper.SetConfigType("toml")
	err := viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

func viperWrite() {
	viper.SetConfigFile("cfg.toml")
	viper.Set("Address", "0.0.0.0:9090") //统一把Key处理成小写 Address->address
	viper.WriteConfig()                  // 将当前配置写入“viper.AddConfigPath()”和“viper.SetConfigName”设置的预定义路径
	viper.SafeWriteConfig()
	viper.WriteConfigAs("/path/cfg/.config")
	viper.SafeWriteConfigAs("/path/cfg/.config") // 因为该配置文件写入过，所以会报错

	err := viper.WriteConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %+v", err))
	}
}
