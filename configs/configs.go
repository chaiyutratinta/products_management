package configs

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Database struct {
		Username     string `yaml:"user"`
		Password     string `yaml:"pass"`
		DatabaseName string `yaml:"name"`
		Server       string `yaml:"server"`
	} `yaml:"database"`
}

func getConf() Conf {
	var confFile string

	if env := os.Getenv("APP_ENV"); env == "production" {
		confFile = "./prod.yml"
	} else {
		confFile = "./dev.yml"
	}

	f, err := os.Open(confFile)

	if err != nil {
		panic(err)
	}

	cfg := Conf{}
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)

	if err != nil {
		panic(err)
	}

	return cfg
}

var Config Conf = getConf()
