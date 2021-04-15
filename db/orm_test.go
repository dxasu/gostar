package db

import (
	"errors"
	"time"

	log "../util/glog"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

type User struct {
	ID        uint      // 列名为 `id`
	Name      string    // 列名为 `name`
	Birthday  time.Time // 列名为 `birthday`
	CreatedAt time.Time // 列名为 `created_at`
	Age       int
}

type Personnice struct {
	ID        uint   `json:”id”`
	FirstName string `json:”firstname”`
	LastName  string `json:”lastname”`
	NewName   string `json:”newname”`
}

func (u *Personnice) BeforeCreate(tx *gorm.DB) (err error) {
	log.Infoln("person", u)
	return
}

func Test_orm() {
	payData, _ := NewORM("mysqlPay")
	payData.AutoMigrate(Personnice{}, User{})
	count := 0
	person := Personnice{21, "2", "3", "4"}
	ret := payData.Clauses(clause.OnConflict{DoNothing: true}).Omit("FirstName").Create(&person)
	log.Infoln(">>>>>1", person.ID, ret.Error, ret.RowsAffected)
	person2 := Personnice{FirstName: "jack", NewName: "newname", LastName: "lastname"}
	ret = payData.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Select("FirstName", "LastName").Create(&person2)
	log.Infoln(">>>>>2", person2.ID, ret.Error, ret.RowsAffected)
	exmp := Personnice{41, "2", "3", "4"}

	//payData.Where("id = ?", 2).Delete(&Personnice{}).Error
	//payData.Model(&Personnice{}).Count(&count)
	payData.First(exmp, 1)
	payData.Delete(exmp, 1)
	log.Infoln("count", count)

	payData.Model(&User{}).Create(map[string]interface{}{
		"Name": "shenzhen", "Age": 18,
	})

	// batch insert from `[]map[string]interface{}{}`
	payData.Model(&User{}).Create([]map[string]interface{}{
		{"Name": "shenzhen_1", "Age": 19},
		{"Name": "shenzhen_2", "Age": 20},
	})

	ret = payData.Take(&person) //// SELECT * FROM person LIMIT 1;
	log.Infof("take: %+v %+v %d\n", person, ret.Error, ret.RowsAffected)
	log.Infoln("errors.Is(result.Error, gorm.ErrRecordNotFound)", errors.Is(ret.Error, gorm.ErrRecordNotFound))
	payData.Limit(1).Find(&person)
	log.Infof("Find: %+v %+v %d\n", person, ret.Error, ret.RowsAffected)

	stmt := payData.Session(&gorm.Session{DryRun: true}).First(&person, 1).Statement
	//=> SELECT * FROM `users` WHERE `id` = $1 ORDER BY `id`
	//=> []interface{}{1}
	log.Infof("payData.Session sql:%s  val：%+v\n", stmt.SQL.String(), stmt.Vars)

}
