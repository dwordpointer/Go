package database

import (
	"main/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	connection, err := gorm.Open(mysql.Open("root:@/go_auth"), &gorm.Config{})

	if err != nil {
		panic("Bağlantı kurulamıyor.")
	}

	DB = connection

	connection.AutoMigrate(&models.User{})
}
