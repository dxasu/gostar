package db

import (
	"fmt"
	"strconv"

	"github.com/dxasu/gostar/config"
	log "github.com/dxasu/gostar/util/glog"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//type ORM gorm.DB
type ORM struct {
	gorm.DB
}

// NewDB get one db by dcCfg
func NewORM(cfgKey string, opts ...gorm.Option) (*ORM, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("InitRedis redisCfg:%s err:%+v\n", cfgKey, err.(error))
		}
	}()

	cfg := config.GetCfgMapStr(cfgKey)

	// db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	DSN := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", cfg["user"], cfg["password"], cfg["host"], cfg["database"])

	db, err := gorm.Open(mysql.Open(DSN), opts...)

	if err != nil {
		log.Errorf("Connection to mysql failed. host:%s database:%s err:%+v\n", cfg["host"], cfg["database"], err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("db.DB() failed. host:%s database:%s err:%+v\n", cfg["host"], cfg["database"], err)
		return nil, err
	}

	maxconn, _ := strconv.Atoi(cfg["maxconn"])
	sqlDB.SetMaxOpenConns(maxconn)
	sqlDB.SetMaxIdleConns(10)
	//sqlDB.SetConnMaxLifetime(time.Hour)
	log.Infof("Connection to mysql success. host:%s database:%s\n", cfg["host"], cfg["database"])
	return &ORM{*db}, nil
}
