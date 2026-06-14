package main

import (
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
	initializers.InitDB()
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
