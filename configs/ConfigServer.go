package configs

// Server ...
type Server struct {
	Port	string "toml:\"port\""
	DB 	*Database "toml:\"database\""
}

// NewConfig ... 
func NewConfig() *Server {
	return &Server{
		Port: "8080",
		DB: &Database{
			Port: "5555",
			DBURL: "host=localhost user=abdr password=qwerty123 dbname=abdr sslmode=disable",
		},
	}
}