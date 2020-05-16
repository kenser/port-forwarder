package model

import (
	"github.com/cloverzrg/go-portforward/db"
	"github.com/cloverzrg/go-portforward/model/forwarddao"
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
