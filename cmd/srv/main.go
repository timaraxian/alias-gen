package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/timaraxian/hotel-gen/pkg/application"
	"github.com/timaraxian/hotel-gen/pkg/helpers/alerts"
)

func main() {
	config := application.Config{}
	if _, err := toml.DecodeFile(os.Getenv("HOTELGEN_CONFIG"), &config); err != nil {
		log.Printf("Failed to open config file: %s\n", err)
		os.Exit(1)
	}

	//application.DBFreshService,
	app, err := application.Mount(config, []application.Service{
		application.DBService,
	})
	if err != nil {
		alerts.AlertError(err, "Failed to mount application")
		os.Exit(3)
	}

	srv := &http.Server{
		ReadTimeout:  32 * time.Minute,
		WriteTimeout: 32 * time.Minute,
		Addr:         config.ListenAddr,
		Handler:      app.Routes(),
	}

	alerts.AlertError(nil, "Starting server at address %s", srv.Addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			alerts.AlertError(err, "Server stopped")
			os.Exit(4)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)

	<-sigChan

	if err := srv.Close(); err != nil {
		alerts.AlertError(err, "Failed stopping server")
	} else {
		alerts.AlertError(nil, "Server closed")
	}
}
