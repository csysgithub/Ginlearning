package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string `json:"name" binding:"required"`
	Email string `json:" email" binding:"required"`
}
type UpdateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	r := gin.Default()
	dsn := "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	//建表
	err = db.AutoMigrate(&User{})

	if err != nil {
		panic("failed to connect database")
	}
	r.Run(":8888")
}
