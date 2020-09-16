package application

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/timaraxian/hotel-gen/pkg/database"
	"github.com/timaraxian/hotel-gen/pkg/errors"
)

type App struct {
	Config Config
	DBAL   *database.DBAL
}

type apiResponse struct {
	OK   bool        `json:"ok"`
	Data interface{} `json:"data"`
}

type Service func(app *App) (err error)

type Config struct {
	Env string

	DB database.Config

	ListenAddr string
}

func Mount(config Config, services []Service) (a *App, err error) {
	a = &App{Config: config}

	for _, s := range services {
		if err = s(a); err != nil {
			return a, err
		}
	}

	return a, err
}

func (app *App) decodeRequest(r *http.Request, args interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		return errors.HttpBadRequestArgs
	}
	return nil
}

func (app *App) respondApi(w http.ResponseWriter, r *http.Request, data interface{}, err error) {
	if err != nil {
		if appErr, ok := err.(errors.Error); !ok {
			app.serverErr(w, r, err)
		} else {
			data = appErr.Code()
		}
	}

	if resp, err := json.Marshal(apiResponse{err == nil, data}); err != nil {
		app.serverErr(w, r, err)
	} else {
		w.Write(resp)
	}
}

/*
func (app *App) executeView(w http.ResponseWriter, r *http.Request, path string, data interface{}) {
	if err := app.Views.Execute(w, path, data); err != nil {
		app.serverErr(w, r, err)
	}
}
*/

func (app *App) serverErr(w http.ResponseWriter, r *http.Request, err error) {
	log.Output(2, fmt.Sprintf("ERROR: %s\n%s", err.Error(), debug.Stack()))

	if r.Header.Get("Content-Type") == "application/json" {
		w.Write([]byte(`{"ok": false, "data": "Unexpected"}`))
	} else {
		w.WriteHeader(500)
		w.Write([]byte("Unexpected"))
	}
}
