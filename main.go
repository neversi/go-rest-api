package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/configs"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/server"
)

var (
	configPath string
)

func init() { 
	flag.StringVar(&configPath, "config", "configs/Configs.toml", "configure program with the options in the passed file")
}

func main() {
	flag.Parse()
	fmt.Println("Hello lets start our server...")
	os.Setenv("TokenPass", "abdr_go_to_env")

	conf := configs.NewConfig()
	toml.DecodeFile(configPath, &conf)
	
	api, err := server.New(conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := api.Start(); err != nil {
		fmt.Println(err)
	}

}
