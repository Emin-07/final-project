package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	router.Handle("/*", http.FileServer(http.Dir("web")))
	router.Get("/api/nextdate", app.getNextDate)

	router.Post("/api/task", app.addTaskHandler)
	router.Get("/api/task", app.taskHandler)
	router.Put("/api/task", app.changeTaskHandler)

	router.Get("/api/tasks", app.tasksHandler)

	return router
}
