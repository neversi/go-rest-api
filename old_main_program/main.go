package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	"github.com/neversi/go-rest-api/app/server"
)

var (
	configPath 	string
	logPath		string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config.json", "path to config file")
}

func main() {
	flag.Parse()

	config := server.NewConfig()
	file, err := ioutil.ReadFile(configPath)
	err = json.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}
	s, err := server.New(config)
	if err != nil {
		panic(err)
	}


	if err := s.Start(); err != nil {
		panic(err)
	}
	
}