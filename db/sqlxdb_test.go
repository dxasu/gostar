package db

import (
	"errors"

	log "github.com/dxasu/gostar/util/glog"

	"github.com/go-redis/redis"
)

// DB.NamedExec方法用来绑定SQL语句与结构体或map中的同名字段
func (this *SQLX) InsertUserDemo() (err error) {
	sqlStr := "INSERT INTO user (name,age) VALUES (:name,:age)"
	_, err = this.NamedExec(sqlStr,
		map[string]interface{}{
			"name": "七米",
			"age":  28,
		})
	return
}

// 插入数据
func (this *SQLX) InsertRowDemo() {
	sqlStr := "insert into user(name, age) values (?,?)"
	ret, err := this.Exec(sqlStr, "china", 19)
	if err != nil {
		log.Infof("insert failed, err:%v\n", err)
		return
	}
	theID, err := ret.LastInsertId() // 新插入数据的id
	if err != nil {
		log.Infof("get lastinsert ID failed, err:%v\n", err)
		return
	}
	log.Infof("insert success, the id is %d.\n", theID)
}

// 更新数据
func (this *SQLX) UpdateRowDemo() {
	sqlStr := "update user set age=? where id = ?"
	ret, err := this.Exec(sqlStr, 39, 6)
	if err != nil {
		log.Infof("update failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		log.Infof("get RowsAffected failed, err:%v\n", err)
		return
	}
	log.Infof("update success, affected rows:%d\n", n)
}

// 删除数据
func (this *SQLX) DeleteRowDemo() {
	sqlStr := "delete from user where id = ?"
	ret, err := this.Exec(sqlStr, 6)
	if err != nil {
		log.Infof("delete failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		log.Infof("get RowsAffected failed, err:%v\n", err)
		return
	}
	log.Infof("delete success, affected rows:%d\n", n)
}

// 事务操作
func (this *SQLX) TransactionDemo2() (err error) {
	tx, err := this.Beginx() // 开启事务
	if err != nil {
		log.Infof("begin trans failed, err:%v\n", err)
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			log.Infoln("rollback")
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
			log.Infoln("commit")
		}
	}()

	sqlStr1 := "Update user set age=20 where id=?"

	rs, err := tx.Exec(sqlStr1, 1)
	if err != nil {
		return err
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("exec sqlStr1 failed")
	}
	sqlStr2 := "Update user set age=50 where i=?"
	rs, err = tx.Exec(sqlStr2, 5)
	if err != nil {
		return err
	}
	n, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("exec sqlStr1 failed")
	}
	return err
}

func Test_Get(client *redis.Client) {
	pong, err := client.Ping().Result()
	log.Infoln(pong, err)

	err = client.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val2, err := client.Get("key2").Result()
	if err == redis.Nil {
		log.Infoln("key2 does not exists")
	} else if err != nil {
		panic(err)
	} else {
		log.Infoln("key2", val2)
	}
}
