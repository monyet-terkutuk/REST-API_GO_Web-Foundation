package main

import (
	"fmt"
	"go_api_foundation/user"
	"log"

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

	userRepository := user.NewRepository(db)
	user := user.User{
		Name: "Test Save",
	}

	userRepository.Save(user)
	fmt.Println("User di tambahkan")

}
