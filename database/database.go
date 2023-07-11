package database

import (
	"log"
	"os"
	"rest-api/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBinstance struct {
	DB *gorm.DB
}

var Database DBinstance

func ConnectDB() {
	db, err := gorm.Open(sqlite.Open("api.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect with db!\n", err.Error())
		os.Exit(2)
	}

	log.Println("Connected to db successfully")
	db.Logger = logger.Default.LogMode(logger.Info)

	db.AutoMigrate(&model.User{}, &model.Product{}, &model.Order{})

	Database = DBinstance{DB: db}
}
