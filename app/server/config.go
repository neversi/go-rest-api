package server

import "github.com/neversi/go-rest-api/app/db"

//Config struct contains the configure metadata of the server
type Config struct {
	BindPort 	string "json:\"bind_port\""
	LogFile 	string "json:\"log_file\""
	DBConf 		*db.Config
	DBPath		string "json:\"db_url\""
}

// NewConfig function creates setted configuration
func NewConfig() *Config {

	return &Config{
		BindPort: ":8080",
		LogFile: "log.txt",
		DBConf: db.NewConfig(),
		DBPath: "",
	}
}