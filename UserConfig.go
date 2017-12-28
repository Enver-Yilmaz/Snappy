package main

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

type UserConfig struct {
	AllucKey string
	RealDebridKey string
	TmdbKey string
}

func ParseConfigFile(path string)(config UserConfig){


	fileContent, err := ioutil.ReadFile(path + "/config.cfg")

	if err != nil {
		panic(err)
	}

	if _, err := toml.Decode(string(fileContent), &config); err != nil {
		panic(err)
	}

	return
}
