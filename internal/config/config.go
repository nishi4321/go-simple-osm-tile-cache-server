package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	TILE struct {
		ORIGIN string `yaml:"ORIGIN"`
		VALUES []struct {
			Key   string `yaml:"KEY"`
			Value string `yaml:"VALUE"`
		} `yaml:"VALUES"`
	} `yaml:"TILE"`
}

var config Config

func init() {
	buf, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Fatalln(err)
	}
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		log.Fatalln(err)
	}
}

func Get() Config {
	return config
}
