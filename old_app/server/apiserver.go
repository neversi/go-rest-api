package server

import (
	"log"
	"net/http"
	"os"

	"github.com/neversi/go-rest-api/app/db"
	// "github.com/neversi/go-rest-api/app/objects"
)

//APIServer struct with configuration
type APIServer struct {
	config *Config
	logger *log.Logger
	dbStore *db.DB

}

//New func creates new instance of APIServer
func New(config *Config) (*APIServer, error) {
	file, err := os.OpenFile("log.txt", os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0666)
	if (err != nil)	{
		return nil, err
	}

	myDb := db.New(config.DBConf)

	myDb.SetURL(config.DBPath)

	err = myDb.Open()



	if err != nil {
		return nil, err
	}
	
	return &APIServer{
		config: config,
		logger: log.New(file,"LOG: ", log.Ldate | log.Ltime | log.Lshortfile),
		dbStore: myDb,
	}, nil
}

//Start Running instance of APIServer, returns error if port is busy or some problems occured
func (s *APIServer) Start() error {
	
	s.logger.Println("Server started...")
	defer s.logger.Println("Server terminated...")

	// u := &objects.User{
	// 	Password: "qweqwe123",
	// 	Login: "dorm",
	// 	Email: "here@mail.ru",
	// 	FirstName: "Abdarrakhman",
	// 	SurName: "Akhmetgali",
	// }
	// http.HandleFunc()

	return http.ListenAndServe(s.config.BindPort, nil)
}

