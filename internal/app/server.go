package app

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func (app *App) NewServer() *http.Server {
	return &http.Server{
		Addr:         app.Cfg.Port,
		Handler:      app.routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
}

func (app *App) routes() http.Handler {
	router := chi.NewRouter()
	app.schedulerHandler.RegisterRoutes(router, app.Cfg.Password)
	return router
}
