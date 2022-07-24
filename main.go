package main

import (
	"gofizzbuzz/api"
	"gofizzbuzz/config"
	"strconv"

	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"

	env "github.com/caarlos0/env/v6"
)

func main() {
	// get config
	cfg := config.Data{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid config")
	}

	// set log level
	lvl, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid log level")
	}
	zerolog.SetGlobalLevel(lvl)

	log.Info().Msgf("%+v\n", cfg)

	persistence := true
	app := http.NewServer(persistence, cfg.RequestDataPath)

	log.Info().Msg("server start")
	err = app.Listen(":" + strconv.Itoa(cfg.ServerPort))
	if err != nil {
		log.Log().Err(err).Msg("could not start the server")
	}
}
