package main

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"github.com/nsf/termbox-go"
	"log"
)

//TODO: add a default file, oh but snappy wouldn't work anyway...

type UserConfig struct {
	AllucKey string
	RealDebridKey string
	TmdbKey string
	AllucSearchLength int
}

//parses the user's config file and loads the result into a UserConfig struct
func ParseConfigFile(path string)(config UserConfig){


	fileContent, err := ioutil.ReadFile(path + "/config.cfg")

	if err != nil {
		termbox.Close()
		log.Fatal(err)
	}

	if _, err := toml.Decode(string(fileContent), &config); err != nil {
		termbox.Close()
		log.Fatal(err)
	}

	return
}
