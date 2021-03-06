package common

import (
	"fmt"
	"lingjiao0710/ginEssential/model"

	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	driverName := "mysql"
	host := "localhost"
	port := "3306"
	database := "ginessential"
	username := "root"
	password := "root1234"
	charset := "utf8"

	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)

	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("连接数据库失败, error: " + err.Error())
	}

	//自动创建数据表
	db.AutoMigrate(&model.User{})

	DB = db
	return db

}

func GetDB() *gorm.DB {
	return DB
}
