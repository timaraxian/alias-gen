package application

import (
	"net/http"
)

func (app *App) Routes() http.Handler {
	mux := http.NewServeMux()

	//	getMdl := middlewareGroup(app.getOnlyWebMdl)

	apiMdl := middlewareGroup(
		app.corsMdl,
		app.postOnlyApiMdl,
		app.jsonOnlyApiMdl,
		app.setApiHeadersMdl,
	)

	// -----------------------------------------------------------------------------
	// Routes
	// -----------------------------------------------------------------------------
	mux.Handle("/wordCreate", apiMdl(http.HandlerFunc(app.WordCreate)))

	return middlewareGroup(
		app.catchPanicMdl,
		app.logMdl,
		app.secureHeadersMdl,
	)(mux)
}
