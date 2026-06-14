package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	router.Handle("/*", http.FileServer(http.Dir("web")))
	router.Get("/api/nextdate", app.getNextDate)

	router.Post("/api/task", auth(app.addTaskHandler))
	router.Get("/api/task", auth(app.taskHandler))
	router.Put("/api/task", auth(app.changeTaskHandler))
	router.Delete("/api/task", auth(app.deleteTaskHandler))

	router.Get("/api/tasks", auth(app.tasksHandler))

	router.Post("/api/task/done", auth(app.completeTaskHandler))

	router.Post("/api/signin", app.signInHandler)

	return router
}
