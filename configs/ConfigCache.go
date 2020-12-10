package configs


// Cache ...
type Cache struct {
	Host 		string	"toml:\"host\""
	Port 		string	"toml:\"port\""
	DB		int 	"toml:\"db\""
	ExpDuration 	int64	"toml:\"exp\""
}