package app

import (
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

	rep := repository.NewPostgresRepository(pg)

	// service
	serv := service.NewMyService(rep)

	// handler
	h := handler.NewHandler(serv)

	router := mux.NewRouter()
	router.HandleFunc("/users", h.CreateUserSegments).Methods("POST")
	router.HandleFunc("/segments", h.CreateSegment).Methods("POST")
	router.HandleFunc("/segments", h.DeleteSegment).Methods("DELETE")
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	http.ListenAndServe(net.JoinHostPort("", cfg.Port), router)
}
