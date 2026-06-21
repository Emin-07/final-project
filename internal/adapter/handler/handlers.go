package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Emin-07/final-project/internal/core/domain"
	"github.com/Emin-07/final-project/internal/core/service"
	"github.com/Emin-07/final-project/internal/pkg/hash"
)

func (sh *SchedulerHandler) getNextDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	if date == "" {
		JsonError(w, "no date passed in parameters", http.StatusBadRequest)
		return
	}
	if repeat == "" {
		JsonError(w, "no repeat passed in parameters", http.StatusBadRequest)
		return
	}

	now, err := service.ValidateStringTime(nowStr)
	if err != nil {
		JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := sh.service.NextDate(now, date, repeat)
	if err != nil {
		JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

func (sh *SchedulerHandler) addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task domain.Task

	data, err := io.ReadAll(r.Body)
	if err != nil {
		JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(data, &task)
	if err != nil {
		JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		JsonError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}
	err = sh.service.CheckDate(&task)
	if err != nil {
		JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := sh.service.AddTask(r.Context(), &task)
	if err != nil {
		JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJson(w, map[string]int64{"id": id}, http.StatusCreated)
}

func (sh *SchedulerHandler) getTask(r *http.Request) (*domain.Task, error) {
	id := r.URL.Query().Get("id")
	if id == "" {
		return nil, noIdProvidedErr
	}
	task, err := sh.service.GetTask(r.Context(), id)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (sh *SchedulerHandler) listTasks(r *http.Request) ([]*domain.Task, error) {
	searchParameter := r.URL.Query().Get("search")
	limitParameter := r.URL.Query().Get("limit")

	return sh.service.GetTasks(r.Context(), limitParameter, strings.TrimSpace(searchParameter))
}

func (sh *SchedulerHandler) tasksHandler(w http.ResponseWriter, r *http.Request) {
	task, err := sh.getTask(r)
	if err != nil {
		if errors.Is(err, noIdProvidedErr) {
			tasks, err := sh.listTasks(r)
			if err != nil {
				JsonError(w, err.Error(), http.StatusBadRequest)
				return
			}
			WriteJson(w, TasksResp{
				Tasks: tasks,
			}, http.StatusOK)
			return
		} else if errors.Is(err, domain.ErrNoRecord) {
			JsonError(w, fmt.Sprintf("нет задания с переданным id: %v", err.Error()), http.StatusNotFound)
			return
		}
		JsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	WriteJson(w, task, http.StatusOK)
}

func (sh *SchedulerHandler) changeTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task domain.Task

	data, err := io.ReadAll(r.Body)
	if err != nil {
		JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(data, &task)
	if err != nil {
		JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		JsonError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	err = sh.service.CheckDate(&task)
	if err != nil {
		JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = sh.service.ChangeTask(r.Context(), &task)
	if err != nil {
		JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	WriteJson(w, map[string]string{}, http.StatusOK)
}

type TasksResp struct {
	Tasks []*domain.Task `json:"tasks"`
}

func (sh *SchedulerHandler) completeTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := sh.getTask(r)
	if err != nil {
		if errors.Is(err, domain.ErrNoRecord) {
			JsonError(w, fmt.Sprintf("нет задания с переданным id: %v", err.Error()), http.StatusNotFound)
			return
		}
		JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = sh.service.CompleteTask(r.Context(), task)
	if err != nil {
		JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJson(w, map[string]string{}, http.StatusOK)
}

func (sh *SchedulerHandler) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		JsonError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}
	err := sh.service.DeleteTask(r.Context(), id)
	if err != nil {
		JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJson(w, map[string]string{}, http.StatusOK)
}

func (sh *SchedulerHandler) signInHandler(w http.ResponseWriter, r *http.Request) {
	userPassword := map[string]string{}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(data, &userPassword)
	if err != nil {
		JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	passHashed, err := hash.HashPassword(os.Getenv("TODO_PASSWORD"))
	if err != nil {
		JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !hash.CheckPasswordHash(userPassword["password"], passHashed) {
		JsonError(w, "Неверный пароль", http.StatusBadRequest)
		return
	}

	tokenStr, err := service.StringToken(passHashed)
	if err != nil {
		JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})
	//fmt.Printf("{ \"token\": %v }\n", tokenStr) // for dev
	WriteJson(w, map[string]string{"token": tokenStr}, http.StatusOK)
}
