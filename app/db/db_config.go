package db


// Config for DB
type Config struct {
	dbURL string
}

// NewConfig instantiates config of the DB
func NewConfig() *Config {
	return &Config{
		dbURL: "",
	}
}
// SetURL sets url of DB
func (db *DB) SetURL(str string) {
	db.config.dbURL = str
}