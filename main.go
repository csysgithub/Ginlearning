package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type User struct {
	gorm.Model
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
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

	r.GET("/users", func(c *gin.Context) {
		var users []User
		if err := db.Find(&users).Error; err != nil {
			log.Println("Failed to find users:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to find users",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": users,
		})
	})

	r.POST("/users", func(c *gin.Context) {
		var input User
		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err := db.Create(&input).Error; err != nil {
			log.Println("Failed to create user:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create user",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": input,
		})
	})

	r.PUT("/users/:id", func(c *gin.Context) {
		var input UpdateUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid input data",
			})
			log.Println("Invalid input data:", err)
			return
		}

		id := c.Param("id")
		var user User
		if err := db.First(&user, id).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			log.Println("User not found:", err)
			return
		}

		if input.Name != "" {
			user.Name = input.Name
		}
		if input.Email != "" {
			user.Email = input.Email
		}
		if err := db.Save(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update user",
			})
			log.Println("Failed to update user:", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": user,
		})
	})

	r.DELETE("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var user User
		if err := db.First(&user, id).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "user not found",
			})
			log.Println("user not found:", err)
			return
		}

		if err := db.Delete(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to delete user",
			})
			log.Println("Failed to delete user:", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": "User deleted successfully",
		})
	})

	r.Run(":8888")
}
