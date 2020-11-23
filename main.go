package main

import (
	"flag"
	"fmt"

	// "gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/database"
	// "gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/models"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/server"
)

var configPath string

func init() { 
	flag.StringVar(&configPath, "configPath", "config.json", "configure file which contains certain details")
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

	// db := database.New()
	// db.OpenDataBase()

	// user := &models.User{
	// 	Login: "abadr",
	// 	Password: "ABDRABDR",
	// 	FirstName: "Abdarrakhman",
	// 	SurName: "Akhmetgali",
	// 	Email: "Abdarrakhman@gmail.com",
	// }
	// _ = user

	// users, _ := db.User().Read(nil)
	// _ = users
}
