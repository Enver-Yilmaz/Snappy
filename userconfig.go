package main

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

//TODO: add a default file, oh but snappy wouldn't work anyway...

type UserConfig struct {
	AllucKey string
	RealDebridKey string
	TmdbKey string
	AllucSearchLength int
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
