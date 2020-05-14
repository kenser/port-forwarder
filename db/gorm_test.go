package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"testing"
)

type Data struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func TestGorm(t *testing.T) {
	db, err := gorm.Open("sqlite3", "./sqlite3.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	data := Data{
		Id:   1,
		Name: "123",
	}
	if !db.HasTable(&data) {
		err = db.CreateTable(&data).Error
		if err != nil {
			t.Fatal(err)
		}
	}
	err = db.Create(&data).Error
	if err != nil {
		t.Fatal(err)
	}

}
