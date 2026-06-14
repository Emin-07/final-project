package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/Emin-07/final-project/internal/initializers"
	"github.com/Emin-07/final-project/internal/models"
	_ "modernc.org/sqlite"
)

type application struct {
	schedule *models.SchedulerModel
}

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	dbfile := os.Getenv("TODO_DBFILE")

	db, err := createTablesAndOpenDB(dbfile)
	if err != nil {
		log.Fatalf("couldn't connect to database: %v", err)
	}
	defer db.Close()

	app := &application{
		schedule: &models.SchedulerModel{DB: db},
	}

	if err = app.Run(); err != nil {
		log.Fatal(err)
	}
}

func createTablesAndOpenDB(dsn string) (*sql.DB, error) {
	// dsn == db name if it's sqlite you're using
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Adding table into db if it's empty
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	isEmpty, err := models.HasNoTables(ctx, db, `SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%'`)
	if err != nil {
		return nil, err
	}

	if isEmpty {
		err = models.InitScheduleTable(ctx, db)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
