package config

import (
	"flag"

	log "github.com/dxasu/gostar/util/glog"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	flag.String("path", "", "Input path *.toml")
	flag.Set("log_dir", "./log")
	flag.Set("alsologtostderr", "true")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	log.Infof("flags %+v", viper.AllSettings())
}
