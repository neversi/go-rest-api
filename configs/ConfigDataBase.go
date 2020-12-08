package configs

import (
	"github.com/!burnt!sushi/toml"
)

// Database ...
type Database struct {
	Port 	string "toml:\"port\""
	DBURL	string "toml:\"database_url\""
}

// NewConfigDB ...
func NewConfigDB() *Database {
	
	return nil
}