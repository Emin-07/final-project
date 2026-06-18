package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"

	"github.com/Emin-07/final-project/internal/adapter/repository/sqlite/scheduler"

	"github.com/Emin-07/final-project/internal/adapter/handler"

	"github.com/Emin-07/final-project/internal/core/service"

	"github.com/Emin-07/final-project/internal/app"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}
}

func main() {
	dbfile := os.Getenv("TODO_DBFILE")

	db, err := OpenDB(dbfile)
	if err != nil {
		log.Fatalf("couldn't connect to database: %v", err)
	}
	defer db.Close()
	newRepo := scheduler.NewSchedulerRepo(db)
	err = newRepo.InitDataIntoDb()
	if err != nil {
		log.Fatalf("couldn't initialize tables to database: %v", err)
	}
	newService := service.NewSchedulerService(newRepo)
	newHandler := handler.NewSchedulerHandler(newService)

	application := app.NewApp(app.WithHandler(newHandler))

	srv := application.NewServer()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	serverErr := make(chan error, 1)
	go func() {
		fmt.Printf("Starting server at http://localhost%s...\n", application.Cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		fmt.Printf("listen and serve: %v\n", err)
	case sig := <-shutdown:
		fmt.Printf("получен сигнал %s, выключаемся", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("shutdown: %w", err)
		}
	}
	return
}

func OpenDB(dsn string) (*sql.DB, error) {
	// dsn == db name if it's sqlite you're using
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
