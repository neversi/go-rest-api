package configs

// Server ...
type Server struct {
	Port	string "toml:\"port\""
	DB 	Database "toml:\"database\""
	

}