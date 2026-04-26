package database

import (
	"github.com/NishLy/go-fiber-boilerplate/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(config *config.Config) {
	dsn := "host=" + config.DBHOST + " user=" + config.DBUSER + " password=" + config.DBPASS + " dbname=" + config.DBNAME + " port=" + config.DBPORT + " sslmode=disable TimeZone=UTC"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	DB = db
}
