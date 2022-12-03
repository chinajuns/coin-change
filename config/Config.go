package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type configs struct {
	APP     interface{} `yaml:"APP"`
	MYSQL   interface{} `yaml:"MYSQL"`
	REDIS   interface{} `yaml:"REDIS"`
	EMAIL   interface{} `yaml:"EMAIL"`
	MONGODB interface{} `yaml:"MONGODB"`
}

// Config
// 解析配置文件
func Config() configs {
	var c configs
	file, err := ioutil.ReadFile("./config/config.yaml")
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = yaml.Unmarshal(file, &c)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return c
}
