package common

import (
	"example.com/ginessential/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	dsn := "DevAuth:Dev127336@tcp(localhost:3306)/ginessential?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database, err: " + err.Error())
	}
	db.AutoMigrate(&model.User{})

	DB = db

	return db
}

func GetDB() *gorm.DB {
	return DB
}
