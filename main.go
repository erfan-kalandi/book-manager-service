package main

import (
	"library/APIHandlers"
	"library/Authenticate"
	"library/config"
	"library/db"
	"net/http"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	var config config.Config

	err := cleanenv.ReadEnv(&config)
	if err != nil {
		logger.WithError(err).Panicln("Can not read config file")
	}

	InitDB, err := db.InitDB(config)
	if err != nil {
		logger.WithError(err).Fatalln("Can not Create datatbase")
	}
	logger.Infoln("Connected to datatbase successfully")

	err = InitDB.CreateSchemas()
	if err != nil {
		logger.WithError(err).Fatalln("Can not migrate tables and models")
	}
	logger.Infoln("automigrate successfully")

	auth, err := Authenticate.NewAuth(InitDB, 10, logger)
	if err != nil {
		logger.WithError(err).Fatal("can not create the authenticate instance")
	}

	bookManagerServer := APIHandlers.Server{
		Authenticate: auth,
		DB:           InitDB,
		Logger:       logger,
	}

	http.HandleFunc("/api/v1/auth/signup", bookManagerServer.HandleSignupAPI)
	http.HandleFunc("/api/v1/auth/login", bookManagerServer.HandleLogin)
	http.HandleFunc("/api/v1/books", bookManagerServer.HandleCreateAndGetAllBook)
	http.HandleFunc("/api/v1/books/", bookManagerServer.HandleUpdateAndDEleteAndGetBook)
	http.HandleFunc("/profile", bookManagerServer.HandleProfile)
	logger.WithError(http.ListenAndServe(":8080", nil)).Fatalln("can not setup the server")

}
