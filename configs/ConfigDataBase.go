package configs


// Database ...
type Database struct {
	Port 	string "toml:\"port\""
	DBURL	string "toml:\"database_url\""
}
