package application

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/timaraxian/alias-gen/pkg/errors"
)

type middleware func(next http.Handler) http.Handler

func middlewareGroup(mdlList ...middleware) middleware {
	return func(handler http.Handler) http.Handler {
		for i := len(mdlList); i > 0; i -= 1 {
			handler = mdlList[i-1](handler)
		}
		return handler
	}
}

func (app *App) secureHeadersMdl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (app *App) logMdl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("INFO: %s - %s %s %s %s", r.RemoteAddr, r.Proto, r.Host, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *App) catchPanicMdl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErr(w, r, errors.Unexpected.WithErr(fmt.Errorf("%s", err)))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *App) getOnlyWebMdl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *App) postOnlyApiMdl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			app.respondApi(w, r, nil, errors.HttpBadMethod)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *App) jsonOnlyApiMdl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
			app.respondApi(w, r, false, errors.HttpBadContentType)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *App) setApiHeadersMdl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("ViewCache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		next.ServeHTTP(w, r)
	})
}

func getOrigin(r *http.Request) (origin string) {
	if origin = r.Header.Get("Origin"); len(origin) > 0 {
		return origin
	}

	if origin = r.Header.Get("Referer"); len(origin) == 0 {
		return ""
	}

	if u, err := url.Parse(origin); err != nil {
		return ""
	} else {
		return u.Scheme + "://" + u.Host
	}
}
func (app *App) corsMdl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := getOrigin(r)

		w.Header().Set("Access-Control-Max-Age", "3600")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Origin, Referer")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
