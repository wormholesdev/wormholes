package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mohitsinghs/wormholes/app"
	"github.com/mohitsinghs/wormholes/auth"
	"github.com/mohitsinghs/wormholes/config"
	"github.com/mohitsinghs/wormholes/factory"
	"github.com/mohitsinghs/wormholes/links"
)

var port int
var cfgFile string

func main() {
	flag.IntVar(&port, "port", 3000, "Port to run")
	flag.StringVar(&cfgFile, "config", "", "Path to non-default config")
	conf, err := config.LoadDefault()
	if err != nil {
		log.Printf("Failed to read config : %v", err)
	}
	config.Merge("WH", conf)
	flag.Parse()

	pgconn := conf.Postgres.Connect()

	authStore := auth.NewStore(pgconn)
	linkStore := links.NewStore(pgconn)

	f := factory.New(&conf.Factory)
	f.TryRestore(linkStore.Ids)

	instance := app.Setup(linkStore, authStore, f)

	go func() {
		app.ShowHeader(port)
		instance.Listen(fmt.Sprintf(":%d", port))
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	if err := instance.Shutdown(); err != nil {
		log.Printf("Error stopping server : %v", err.Error())
	} else {
		log.Println("Server Stopped")
	}

	if err := f.Backup(); err != nil {
		log.Printf("Error during backup : %v", err.Error())
	} else {
		log.Println("Backup Successful")
	}
}
