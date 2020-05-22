package config

import (
	"encoding/json"
	"github.com/cloverzrg/go-portforwarder/logger"
	"io/ioutil"
	"os"
	"reflect"
)

type config struct {
	Nameserver string `json:"nameserver"`
	ENV        string `json:"env"`
	HTTP       struct {
		Listen string `json:"listen"`
	}
}

var Config config

func Parse(path string) (err error){
	err = parse(&Config, path)
	if err != nil {
		return err
	}
	return err
}

func parse(config interface{}, path string) (err error) {
	if reflect.ValueOf(config).Kind() != reflect.Ptr {
		logger.Error("config 请传入指针")
		return
	}
	logger.Info("Reading config from " + path)
	file, err := os.Open(path)
	if err != nil {
		logger.Error(err)
		return err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			logger.Error(err)
		}
	}()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("load configuration file ", path, " successfully")

	return
}
