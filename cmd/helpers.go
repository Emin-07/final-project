package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func getEnv(key string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}
	val := os.Getenv(key)
	if val == "" {
		return "", nil
	}
	return val, nil
}

func initDb(dbFile string) (bool, error) {
	_, err := os.Stat(dbFile)
	if err == nil {
		return false, nil
	}
	_, err = os.Create(dbFile)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (app *application) initTables(ctx context.Context) error {
	err := app.schedule.InitScheduleTable(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) writeJson(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func (app *application) jsonError(w http.ResponseWriter, error string, status int) {
	app.writeJson(w, map[string]string{"error": error}, status)
}
