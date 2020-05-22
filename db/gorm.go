package db

import (
	"github.com/cloverzrg/go-portforwarder/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

func Connect() (err error) {
	DB, err = gorm.Open("sqlite3", "./data/sqlite3.db")
	if err != nil {
		logger.Error(err)
		return err
	}
	return err
}
