package configs

// Server ...
type Server struct {
	Port	string "toml:\"port\""
	DB 	*Database "toml:\"database\""
	Cache	*Cache "toml:\"cache\""
}

// NewConfig ... 
func NewConfig() *Server {
	return &Server{
		Port: "8080",
		DB: &Database{
			Port: "5555",
			DBURL: "host=localhost user=abdr password=qwerty123 dbname=abdr sslmode=disable",
		},
		Cache: &Cache{
			Host: "localhost",
			Port: "5444",
			ExpDuration: 1200,
		},
	}
}