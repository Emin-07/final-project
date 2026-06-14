package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Emin-07/final-project/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func (app *application) getNextDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	if date == "" {
		app.jsonError(w, "no date passed in parameters", http.StatusBadRequest)
		return
	}
	if repeat == "" {
		app.jsonError(w, "no repeat passed in parameters", http.StatusBadRequest)
		return
	}
	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(models.DateFormat, nowStr)
		if err != nil {
			app.jsonError(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	res, err := models.NextDate(now, date, repeat)
	if err != nil {
		app.jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

func checkDate(task *models.Task) error {
	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	if task.Date == "" {
		task.Date = now.Format(models.DateFormat)
	}
	t, err := time.Parse(models.DateFormat, task.Date)
	if err != nil {
		return err
	}

	if now.After(t) {
		if len(task.Repeat) == 0 {
			// если правила повторения нет, то берём сегодняшнее число
			task.Date = now.Format(models.DateFormat)
		} else {
			// в противном случае, берём вычисленную ранее следующую дату
			next, err := models.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
	}

	return nil
}

func (app *application) addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	data, err := io.ReadAll(r.Body)
	if err != nil {
		app.jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(data, &task)
	if err != nil {
		app.jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		app.jsonError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}
	err = checkDate(&task)
	if err != nil {
		app.jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := app.schedule.Add(r.Context(), &task)
	if err != nil {
		app.jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	app.writeJson(w, map[string]int64{"id": id}, http.StatusCreated)
}

func (app *application) getTask(r *http.Request) (*models.Task, error) {
	id := r.URL.Query().Get("id")
	if id == "" {
		return nil, fmt.Errorf("Не указан идентификатор")
	}
	task, err := app.schedule.Get(r.Context(), id)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (app *application) taskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := app.getTask(r)
	if err != nil {
		app.jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	app.writeJson(w, task, http.StatusOK)
}

func (app *application) changeTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	data, err := io.ReadAll(r.Body)
	if err != nil {
		app.jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(data, &task)
	if err != nil {
		app.jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		app.jsonError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	err = checkDate(&task)
	if err != nil {
		app.jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = app.schedule.Update(r.Context(), &task)
	if err != nil {
		app.jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	app.writeJson(w, map[string]string{}, http.StatusOK)
}

type TasksResp struct {
	Tasks []*models.Task `json:"tasks"`
}

func (app *application) tasksHandler(w http.ResponseWriter, r *http.Request) {
	searchParameter := r.URL.Query().Get("search")
	limitParameter := r.URL.Query().Get("limit")

	tasks, err := app.schedule.Tasks(r.Context(), limitParameter, strings.TrimSpace(searchParameter))
	if err != nil {
		app.jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	app.writeJson(w, TasksResp{
		Tasks: tasks,
	}, http.StatusOK)

}

func (app *application) completeTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := app.getTask(r)
	if err != nil {
		app.jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if task.Repeat == "" {
		err = app.schedule.Delete(r.Context(), task.ID)
		if err != nil {
			app.jsonError(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		now := time.Now()
		nextDate, err := models.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			app.jsonError(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = app.schedule.UpdateDate(r.Context(), nextDate, task.ID)
		if err != nil {
			app.jsonError(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	app.writeJson(w, map[string]string{}, http.StatusOK)
}

func (app *application) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		app.jsonError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}
	err := app.schedule.Delete(r.Context(), id)
	if err != nil {
		app.jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	app.writeJson(w, map[string]string{}, http.StatusOK)
}

func (app *application) signInHandler(w http.ResponseWriter, r *http.Request) {
	userPassword := map[string]string{}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		app.jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(data, &userPassword)
	if err != nil {
		app.jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	actualPassword := os.Getenv("TODO_PASSWORD")
	passHashed, err := hashPassword(actualPassword)
	if err != nil {
		app.jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !checkPasswordHash(userPassword["password"], passHashed) {
		app.jsonError(w, "Неверный пароль", http.StatusBadRequest)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		Subject:   passHashed,
	})

	key := os.Getenv("JWT_KEY")

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		app.jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})
	//fmt.Printf("{ \"token\": %v }\n", tokenString) // for dev
	app.writeJson(w, map[string]string{"token": tokenString}, http.StatusOK)
}
