package db

import (
	"github.com/cloverzrg/go-portforward/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

func Connect() (err error) {
	DB, err = gorm.Open("sqlite3", "./db/sqlite3.db")
	if err != nil {
		logger.Error(err)
		return err
	}
	return err
}
