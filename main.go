package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mohitsinghs/wormholes/app"
	"github.com/mohitsinghs/wormholes/config"
	"github.com/mohitsinghs/wormholes/constants"
	"github.com/mohitsinghs/wormholes/factory"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var port int
var cfgFile string

func main() {
	app.ShowNotice()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	flag.IntVar(&port, "port", constants.DEFAULT_PORT, "Port to run")
	flag.StringVar(&cfgFile, "config", "", "Path to non-default config")
	conf, err := config.Load(cfgFile)
	if err != nil {
		log.Error().Err(err).Msg("failed to read config")
	}
	config.Merge(constants.ENV_PREFIX, conf)
	flag.Parse()

	f := factory.New(&conf.Factory)
	instance := app.Setup(conf, f)

	go func() {
		instance.Listen(fmt.Sprintf(":%d", port))
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	if err := instance.Shutdown(); err != nil {
		log.Error().Err(err).Msg("error stopping server")
	} else {
		log.Info().Msg("server stopped")
	}
	f.Backup()
}
