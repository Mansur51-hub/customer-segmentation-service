package app

import (
	"context"
	"github.com/Mansur51-hub/customer-segmentation-service/config"
	_ "github.com/Mansur51-hub/customer-segmentation-service/docs"
	"github.com/Mansur51-hub/customer-segmentation-service/handler"
	"github.com/Mansur51-hub/customer-segmentation-service/pkg/logger"
	"github.com/Mansur51-hub/customer-segmentation-service/pkg/postgres"
	"github.com/Mansur51-hub/customer-segmentation-service/repository"
	"github.com/Mansur51-hub/customer-segmentation-service/service"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"
	"net"
	"net/http"
)

func Run(configPath string) {
	cfg, err := config.NewConfig(configPath)

	if err != nil {
		log.Fatal().Err(err).Msg("config")
	}

	// logger
	logger.SetLevel(cfg.Log.Level)

	// postgres
	log.Info().Msg("connecting to postgres...")
	pg, err := postgres.NewPostgres(cfg.Url)

	if err != nil {
		log.Fatal().Err(err).Msg("postgres connect")
	}

	defer pg.Close()

	repos := repository.NewRepositories(pg)

	// service
	serv := service.NewServices(repos)

	// init ttl service

	serv.TtlService.Exec(context.Background())

	// handler
	h := handler.NewHandler(serv)

	router := mux.NewRouter()
	router.HandleFunc("/users", h.CreateUserSegments).Methods("POST")
	router.HandleFunc("/users", h.GetUserSegments).Methods("GET")
	router.HandleFunc("/segments", h.CreateSegment).Methods("POST")
	router.HandleFunc("/segments", h.DeleteSegment).Methods("DELETE")
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	router.HandleFunc("/operations", h.GetOperations).Methods("POST")
	http.ListenAndServe(net.JoinHostPort("", cfg.Port), router)
}
