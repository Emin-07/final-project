package initializers

import (
	"log"
	"os"
)

func InitDB() {
	dbName := os.Getenv("TODO_DBFILE")
	if dbName == "" {
		dbName = "scheduler.db"
	}
	_, err := os.Stat(dbName)
	if err != nil {
		_, err = os.Create(dbName)
		if err != nil {
			log.Fatal("couldn't create database")
		}
	}
}
