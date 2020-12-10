package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/configs"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DataBase ...
type DataBase struct {
	config *configs.Database
	Pdb *gorm.DB
}


// New creates a new database
func New() *DataBase {
	newDB := DataBase{}
	return &newDB
}
// OpenDataBase ... (do not forget logger + )
func (DB *DataBase) OpenDataBase(conf *configs.Database) error {
	var err error

	DB.config = conf
	dsn := fmt.Sprintf("port=%s %s", DB.config.Port, DB.config.DBURL)
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

	return nil
}

// assertTables asserts that tables exist, otherwise creates them
func (DB *DataBase) assertTables() error {
	var err error
	if err = DB.Pdb.AutoMigrate(&models.User{}, &models.Task{}, &models.Role{}); err != nil {
		return err
	}
	return nil 
}
