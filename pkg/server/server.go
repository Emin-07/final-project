package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func Run() error {
	dbFile := "../../scheduler.db"
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	if install {
		_, err := os.Create("scheduler.db")
		if err != nil {
			fmt.Errorf(err.Error())
		}
	}

	// TODO: read create_schedule_table and put it into string and execute afterwards
	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = ":7540"
	}

	router := chi.NewRouter()
	router.Handle("/*", http.FileServer(http.Dir("web")))

	srv := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// graceful shutdown: ждём сигнал от ОС и аккуратно закрываем сервер
	// signal.Notify пересылает SIGINT/SIGTERM в канал, srv.Shutdown ждёт завершения активных запросов
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	serverErr := make(chan error, 1)
	go func() {
		fmt.Printf("Starting server at http://localhost%s...\n", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return fmt.Errorf("listen and serve: %w", err)
	case sig := <-shutdown:
		log.Printf("получен сигнал %s, выключаемся", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("shutdown: %w", err)
		}
	}

	return nil
}
