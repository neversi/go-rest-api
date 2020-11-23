package database

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DataBase ...
type DataBase struct {
	sync.Mutex
	Ctx context.Context
	Pdb *gorm.DB
	Rdb *redis.Client
	userRep *UserRep
	taskRep *TaskRep
}

// New creates a new database
func New() *DataBase {
	newDB := DataBase{}
	newDB.userRep = newUserRep(&newDB)
	newDB.taskRep = newTaskRep(&newDB)
	return &newDB
}
// OpenDataBase ... (do not forget logger)
func (DB *DataBase) OpenDataBase() error {
	var err error
	
	DB.Ctx = context.Background()
	dsn := "port=5432 host=localhost user=abdr password=qwerty123 dbname=abdr sslmode=disable"
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

// User returns the user_manager
func (DB *DataBase) User() (*UserRep) {
	return DB.userRep
}

// Task returns the task_manager 
func (DB *DataBase) Task() (*TaskRep) {
	return DB.taskRep
}

// EncryptString encrypts the string
func EncryptString(str string) (string, error) {
	hBytes, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hBytes), nil
}

// CompareHash asserts the correctness of the encrypted argument
func CompareHash(hashed string, raw string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(
		[]byte(hashed), 
		[]byte(raw),
		); err != nil {
			return false, err
		}
	
	return true, nil
}
