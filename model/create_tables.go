package model

import (
	"github.com/cloverzrg/go-portforwarder/db"
	"github.com/cloverzrg/go-portforwarder/model/forwarddao"
)

func CreateAllTable() (err error) {
	if !db.DB.HasTable(&forwarddao.Forward{}) {
		err = db.DB.CreateTable(&forwarddao.Forward{}).Error
		if err != nil {
			return err
		}
	}
	return err
}
