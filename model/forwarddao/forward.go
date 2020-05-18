package forwarddao

import (
	"github.com/cloverzrg/go-portforward/db"
	"time"
)

type Forward struct {
	Id               int    `gorm:"primary_key; AUTO_INCREMENT"`
	Network          string `gorm:"not null"`
	ListenAddress    string `gorm:"not null"`
	ListenPort       int    `gorm:"not null"`
	TargetAddress    string `gorm:"not null"`
	TargetPort       int    `gorm:"not null"`
	ConnCount        uint   `gorm:"not null"`
	CurrentConnCount uint   `gorm:"not null"`
	Status           int    `gorm:"not null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time `sql:"index"`
}

func (f Forward) TableName() string {
	return "forwards"
}

func Add(network string, listenAddress string, listenPort int, targetAddress string, targetPort int) (id int, err error) {
	data := Forward{
		Network:       network,
		ListenAddress: listenAddress,
		ListenPort:    listenPort,
		TargetAddress: targetAddress,
		TargetPort:    targetPort,
	}
	err = db.DB.Create(&data).Error
	return data.Id, err
}

func GetById(id int) (data Forward, err error) {
	err = db.DB.Model(&Forward{}).Where("id = ?", id).First(&data).Error
	return data, err
}

func UpdateById(id int, data Forward) (err error) {
	return db.DB.Model(&Forward{}).Where("id = ?", id).Update(data).Error
}

func UpdateByIdMap(id int, m map[string]interface{}) (err error) {
	return db.DB.Model(&Forward{}).Where("id = ?", id).Update(m).Error
}


func DeleteById(id int) (err error) {
	return db.DB.Where("id = ?", id).Delete(&Forward{}).Error
}