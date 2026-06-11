package main

import (
	"context"
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
