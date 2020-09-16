package application

import (
	"database/sql"

	"github.com/timaraxian/hotel-gen/pkg/database"
)

// -----------------------------------------------------------------------------
func DBService(app *App) (err error) {
	app.DBAL, err = database.Bootstrap(app.Config.DB)
	return err
}

func NewTestDBService(conn *sql.DB) Service {
	return func(app *App) (err error) {
		app.DBAL = &database.DBAL{DB: conn}
		return app.DBAL.Fresh()
	}
}

// -----------------------------------------------------------------------------
func DBFreshService(app *App) (err error) {
	return app.DBAL.Fresh()
}

// -----------------------------------------------------------------------------
