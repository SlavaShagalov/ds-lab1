package main

import (
	personsDelivery "github.com/SlavaShagalov/ds-lab1/internal/persons/delivery/http"
	personsRepository "github.com/SlavaShagalov/ds-lab1/internal/persons/repository/pgx"
	"github.com/SlavaShagalov/ds-lab1/internal/pkg/config"
	postgres "github.com/SlavaShagalov/ds-lab1/internal/pkg/storages"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"

	mw "github.com/SlavaShagalov/ds-lab1/internal/middleware"

	pLog "github.com/SlavaShagalov/ds-lab1/internal/pkg/log/prod"
)

// main godoc
//
//	@title						Persons API
//
//	@version					1.0
//	@description				Persons API documentation.
//	@termsOfService				http://127.0.0.1/terms
//
//	@contact.name				Persons API Support
//	@contact.url				http://127.0.0.1/support
//	@contact.email				ppersons-support@vk.com
//
//	@host						127.0.0.1
//	@BasePath					/api/v1
func main() {
	// ===== Configuration =====
	config.SetDefaultPostgresConfig()
	viper.SetConfigName("api")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/configs")
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Failed to read configuration: %v\n", err)
		os.Exit(1)
	}
	log.Printf("Configuration read successfully")

	// ===== Logger =====
	logger := pLog.NewDevelopLogger()
	defer func() {
		err = logger.Sync()
		if err != nil {
			log.Println(err)
		}
	}()
	logger.Info("API service starting...")

	// ===== Data Storage =====
	db, err := postgres.NewStd(logger)
	if err != nil {
		os.Exit(1)
	}
	defer func() {
		db.Close()
		logger.Info("Postgres connection closed")
	}()

	personsRepo := personsRepository.New(db, logger)

	accessLog := mw.NewAccessLog(logger)
	cors := mw.NewCors()

	router := mux.NewRouter()

	// ===== Delivery =====
	personsDelivery.RegisterHandlers(router, personsRepo, logger)

	// ===== Swagger =====
	//router.PathPrefix(constants.ApiPrefix + "/swagger/").Handler(httpSwagger.WrapHandler).Methods(http.MethodGet)

	// ===== Router =====
	server := http.Server{
		Addr:    ":" + viper.GetString(config.ServerPort),
		Handler: accessLog(cors(router)),
	}

	// ===== Start =====
	logger.Info("API service started", zap.String("port", viper.GetString(config.ServerPort)))
	if err = server.ListenAndServe(); err != nil {
		logger.Error("API server stopped", zap.Error(err))
	}
}
