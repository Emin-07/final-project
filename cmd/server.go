package main

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
)

func (app *application) Run() error {

	port := os.Getenv("TODO_PORT")

	srv := &http.Server{
		Addr:         port,
		Handler:      app.routes(),
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
