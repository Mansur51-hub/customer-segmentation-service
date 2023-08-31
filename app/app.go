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
	"github.com/rs/zerolog/log"
	oslog "log"
	"net"
	"net/http"
)

func Run() {
	cfg, err := config.NewConfig()

	if err != nil {
		log.Fatal().Err(err).Msg("config")
	}

	// logger
	logger.SetLevel(cfg.Level)

	// postgres
	log.Info().Msg("connecting to postgres...")
	pg, err := postgres.NewPostgres(cfg.Url)

	if err != nil {
		log.Fatal().Err(err).Msg("postgres connect")
	}

	defer pg.Close()

	log.Info().Msg("initiating repositories...")
	repos := repository.NewRepositories(pg)

	log.Info().Msg("initiating services...")
	// service
	serv := service.NewServices(repos)
	// init ttl service
	serv.TtlService.Exec(context.Background())

	log.Info().Msg("initiating handler...")
	// handler
	h := handler.NewHandler(serv)

	// router
	log.Info().Msg("initiating router")

	router := handler.NewRouter(h)

	log.Info().Str("port", cfg.Port).Msg("starting listen server")
	oslog.Fatal(http.ListenAndServe(net.JoinHostPort("", cfg.Port), router.Router))
}
