package main

import (
	"flag"
	"fmt"

	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/server"
)

var (
	serverConfigPath string
	dbConfigPath 	 string	
)

func init() { 
	flag.StringVar(&serverConfigPath, "serverConfig", "configs/Server.toml", "configure server with the options in the passed file")
	flag.StringVar(&dbConfigPath, "dbConfig", "configs/DataBase.toml", "configure data base with the options in the passed file")
}

func main() {
	flag.Parse()
	fmt.Println("Hello lets start our server...")

	api, err := server.New()
	if err != nil {
		fmt.Println(err)
		return
	}

	api.Start()

}
