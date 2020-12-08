package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/configs"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DataBase ...
type DataBase struct {
	config *configs.Database
	Ctx context.Context
	Pdb *gorm.DB
	Rdb *redis.Client
}


// New creates a new database
func New() *DataBase {
	newDB := DataBase{}
	if err := newDB.OpenDataBase(); err != nil {
		return nil;
	}
	return &newDB
}
// OpenDataBase ... (do not forget logger + )
func (DB *DataBase) OpenDataBase() error {
	var err error
	
	DB.Ctx = context.Background()
	dsn := fmt.Sprintf("%s %s", DB.config.Port, DB.config.DBURL)
	newLogger := logger.New(
		log.New(os.Stdout, "\n", log.LstdFlags),logger.Config{
			SlowThreshold: time.Second,
			LogLevel: logger.Info,
			Colorful: true,
		})
	
	if DB.Pdb, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	}); err != nil {
		return err
	}
	if err = DB.assertTables(); err != nil {
		return err
	}

	DB.Rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:5054",
		Password: "",
		DB: 0,
	})
	
	return nil
}

// assertTables asserts that tables exist, otherwise creates them
func (DB *DataBase) assertTables() error {
	var err error
	if err = DB.Pdb.AutoMigrate(&models.User{}, &models.Task{
	}); err != nil {
		return err
	}
	return nil 
}
