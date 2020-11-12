package db

import (
	"database/sql"
	_ "github.com/lib/pq" // Justifying it)
)

// DB struct
type DB struct {
	config *Config
	db     *sql.DB
	userRep *UserRep
}

// New creates DB
func New(config *Config) *DB {

	return &DB{
		config: config,
	}
}

// Open DB
func (db *DB) Open() error {
	// currentDb, err := sql.Open("postgres", db.config.dbURL)
	currentDb, err := sql.Open("postgres", db.config.dbURL)

	if err != nil {
		return err
	}

	if err := currentDb.Ping(); err != nil {
		return err
	}

	
	db.db = currentDb

	

	return nil

}

// Close closes DB
func (db *DB) Close() {
	db.db.Close()
}

// User here
func (db *DB) User() *UserRep {
	if db.userRep != nil {
		return db.userRep
	}
	db.userRep = &UserRep{
		dbUser: db,
	}
	return db.userRep
}	