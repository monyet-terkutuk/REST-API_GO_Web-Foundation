package main

import (
	"go_api_foundation/handler"
	"go_api_foundation/user"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// bisa buat logging dan helper
	dsn := "root:@tcp(127.0.0.1:3306)/go_foundation?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	// menghubungkan struct ke database
	userRepository := user.NewRepository(db)

	// memaping data dari input user untuk di masukan ke userRepository
	userService := user.NewService(userRepository)

	// mengambil data mentah dari client dan di convert dari json ke userService
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()

	// grouping router
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.UserLogin)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)

	router.Run()
}
