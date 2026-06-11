package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Emin-07/final-project/internal/models"
	_ "modernc.org/sqlite"
)

type application struct {
	schedule *models.SchedulerModel
}

func main() {
	dbfile, err := getEnv("TODO_DBFILE")
	if err != nil {
		fmt.Println(err.Error())
	}
	db, dbCreated, err := openDB(dbfile)
	if err != nil {
		log.Fatal(err)
	}

	app := &application{
		schedule: &models.SchedulerModel{DB: db},
	}
	if dbCreated {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err = app.initTables(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}

	defer db.Close()
	if err = app.Run(); err != nil {
		log.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, bool, error) {
	// dsn == db name if it's sqlite you're using
	dbCreated, err := initDb(dsn)
	if err != nil {
		return nil, false, err
	}
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, false, err
	}
	if err = db.Ping(); err != nil {
		return nil, false, err
	}
	return db, dbCreated, nil
}
