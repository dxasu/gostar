package db

import (
	sql "database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/dxasu/gostar/config"
	log "github.com/dxasu/gostar/util/glog"

	"github.com/jmoiron/sqlx"
)

type SQLX struct {
	sqlx.DB
}

// NewDB get one db by dcCfg
func NewDB(cfgKey string) (*SQLX, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("InitRedis redisCfg:%s err:%+v\n", cfgKey, err)
		}
	}()

	cfg := config.GetCfgMapStr(cfgKey)

	connURI := fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg["user"], cfg["password"], cfg["host"], cfg["database"])

	db, err := sqlx.Open("mysql", connURI)
	if err != nil {
		log.Errorf("Connection to mysql failed. host:%s database:%s err:%+v\n", cfg["host"], cfg["database"], err)
		return nil, err
	}

	maxconn, _ := strconv.Atoi(cfg["maxconn"])
	db.SetMaxOpenConns(maxconn)
	log.Infof("Connection to mysql success. host:%s database:%s\n", cfg["host"], cfg["database"])
	return &SQLX{*db}, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
type K map[string]interface{}

func (maps K) GetKeys() []string {
	var keys []string
	for k := range maps {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// GetKeys
func GetJoin(keys []string, keym, separator string) (valstr string) {
	var keyslice []string
	for _, name := range keys {
		keyslice = append(keyslice, strings.ReplaceAll(keym, "?", name))
	}
	valstr = strings.Join(keyslice, separator)
	return
}

// ToStringMap
func ToStringMap(data K) map[string]string {
	strmap := make(map[string]string)
	for k, v := range data {
		if v != nil {
			strmap[k] = string(v.([]byte))
		}
	}
	return strmap
}

// GetKeys key(int string in default)
func GetGrepKeys(arg K) (valstr string) {
	for name, val := range arg {
		if value, ok := val.(string); ok {
			valstr += fmt.Sprintf("%s='%s' AND ", name, value)
		} else if value, ok := val.(int64); ok {
			valstr += fmt.Sprintf("%s=%d AND ", name, value)
		} else if value, ok := val.(int32); ok {
			valstr += fmt.Sprintf("%s=%d AND ", name, value)
		} else if value, ok := val.(int); ok {
			valstr += fmt.Sprintf("%s=%d AND ", name, value)
		}
	}

	valstr = valstr[:len(valstr)-4]
	return
}

func checkin(key string, keys []string) bool {
	for _, v := range keys {
		if v == key {
			return true
		}
	}
	return false
}

// GetUpdateKeys key(int string in default)
func GetUpdateKeys(arg K, keys ...string) (sqlstr string, vals []interface{}) {
	var namestr string
	sqlstr = "%s where "
	var conds []interface{}
	for name, val := range arg {
		if checkin(name, keys) {
			sqlstr += name + " = ? AND "
			conds = append(conds, val)
		} else {
			namestr += name + " = ? ,"
			vals = append(vals, val)
		}
	}

	vals = append(vals, conds...)
	namestr = namestr[:len(namestr)-1]
	sqlstr = sqlstr[:len(sqlstr)-4]
	sqlstr = fmt.Sprintf(sqlstr, namestr)
	return
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// InsertByMap insert
func (master *SQLX) InsertByMap(table string, data K) (lastID int64, err error) {
	log.Infof("SaveByKey table:%s sqldata:%+v\n", table, data)
	keys := data.GetKeys()
	key := GetJoin(keys, "?", ",")
	val := GetJoin(keys, ":?", ",")
	sqlStr := fmt.Sprintf(`INSERT INTO %s (%s) VALUES(%s)`, table, key, val)
	var ret sql.Result
	ret, err = master.NamedExec(sqlStr, data)
	if err != nil {
		log.Infof("insert failed, sql:%s err:%v\n", sqlStr, err)
		return
	}

	lastID, err = ret.LastInsertId()
	if err != nil {
		log.Infof("get lastinsert ID failed, sql:%s err:%v\n", sqlStr, err)
		return
	}
	log.Infof("insert success, the id is %d.\n", lastID)
	return
}

// SaveByMap save
func (master *SQLX) SaveByMap(table string, data K) (lastID int64, err error) {
	log.Infof("SaveByKey table:%s sqldata:%+v\n", table, data)
	keys := data.GetKeys()
	key := GetJoin(keys, "?", ",")
	val := GetJoin(keys, ":?", ",")
	updatestr := GetJoin(keys, "?=VALUES(?)", ",")
	sqlStr := fmt.Sprintf(`INSERT INTO %s (%s) VALUES(%s) ON DUPLICATE KEY UPDATE %s`, table, key, val, updatestr)
	var ret sql.Result
	ret, err = master.NamedExec(sqlStr, data)
	if err != nil {
		log.Infof("insert failed, sql:%s err:%v\n", sqlStr, err)
		return
	}

	lastID, err = ret.LastInsertId()
	if err != nil {
		log.Infof("get lastinsert ID failed, sql:%s err:%v\n", sqlStr, err)
		return
	}
	log.Infof("insert success, the id is %d.\n", lastID)
	return
}

// UpdateByMap
func (master *SQLX) UpdateByMap(table string, data K, keys ...string) (lastID int64, err error) {
	log.Infof("SaveByKey table:%s sqldata:%+v\n", table, data)
	updatestr, vals := GetUpdateKeys(data, keys...)
	sqlStr := fmt.Sprintf(`UPDATE %s SET %s`, table, updatestr)
	var ret sql.Result
	ret, err = master.Exec(sqlStr, vals...)
	if err != nil {
		log.Infof("insert failed, sql:%s err:%v\n", sqlStr, err)
		return
	}

	lastID, err = ret.LastInsertId()
	if err != nil {
		log.Infof("get lastinsert ID failed, sql:%s err:%v\n", sqlStr, err)
		return
	}
	log.Infof("insert success, the id is %d.\n", lastID)
	return
}

// GetByMap
func (master *SQLX) GetByMap(table string, data K, prop ...string) (map[string]string, error) {
	log.Infof("GetByMap table:%s sqldata:%+v\n", table, data)

	propstr := "*"
	if len(prop) > 0 {
		propstr = strings.Join(prop, ",")
	}

	keystr := GetGrepKeys(data)

	sql := fmt.Sprintf(`SELECT %s FROM %s WHERE %s LIMIT 1`, propstr, table, keystr)

	log.Infof("GetByMap sql:%s\n", sql)

	recvdata := make(K)

	rows, err := master.Queryx(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.MapScan(recvdata)
		break
	}
	if err == nil {
		return ToStringMap(recvdata), nil
	}

	return nil, err
}

// TODO
// func (master *SQLX) GetManyByMap(table string, data K) ([]K, error) {
// 	log.Infof("GetManyByMap table:%s sqldata:%+v\n", table, data)
// 	key, val := GetGrepKeys(data)

// 	sql := fmt.Sprintf(`SELECT %s FROM %s WHERE %s`, key, table, val)

// 	recvdata := make([]K, 0, 4)
// 	err := master.Select(&recvdata, sql)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return recvdata, err
// }
