package db

import (
	"fileserver/internal/auth/model"
	"fileserver/internal/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var conf config.Config

func init() {
	var err error
	conf, err = config.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}
}

func InitDB() *gorm.DB {

	dsn := conf.DbString

	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)

	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(
		&model.User{},
	)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
