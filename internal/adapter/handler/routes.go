package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (sh *SchedulerHandler) RegisterRoutes(router *chi.Mux, passwordForAuth string) http.Handler {

	router.Handle("/*", http.FileServer(http.Dir("./web")))
	router.Get("/api/nextdate", sh.getNextDate)

	router.Post("/api/tasks", Auth(sh.addTaskHandler, passwordForAuth))
	router.Get("/api/tasks/*", Auth(sh.tasksHandler, passwordForAuth))
	router.Put("/api/tasks", Auth(sh.changeTaskHandler, passwordForAuth))
	router.Delete("/api/tasks", Auth(sh.deleteTaskHandler, passwordForAuth))

	router.Post("/api/tasks/done", Auth(sh.completeTaskHandler, passwordForAuth))

	router.Post("/api/signin", sh.signInHandler)

	return router
}
