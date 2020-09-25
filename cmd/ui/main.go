package main

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/timaraxian/alias-gen/pkg/database"
	"github.com/timaraxian/alias-gen/pkg/tui"
)

func main() {
	config := tui.Config{}
	if _, err := toml.DecodeFile(os.Getenv("ALIASGEN_CONFIG"), &config); err != nil {
		log.Printf("Failed to open config file: %s\n", err)
		os.Exit(1)
	}

	dbal, err := database.Bootstrap(config.DB)
	if err != nil{
		panic(err)
	}

	app, err := tui.NewApp(config, []tui.Service{
		func(a *tui.App) (err error) {
			a.DBAL = dbal
			return nil
		},
	})
	if err != nil {
		panic(err)
	}

	if err := app.Loop(); err != nil {
		panic(err)
	}
}
