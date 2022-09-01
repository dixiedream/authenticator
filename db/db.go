package db

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/dixiedream/authenticator/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	p := os.Getenv("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		panic(err)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=true", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), port, os.Getenv("DB_DATABASE"))
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Panic(err.Error())
		panic("failed to connect database")
	}

	// log.Println("Connection Opened to Database")
	DB.AutoMigrate(&model.Server{}, &model.Session{})
	// log.Println("Database Migrated")
}
