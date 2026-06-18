package app

import (
	"os"

	"github.com/Emin-07/final-project/internal/adapter/handler"
)

type Config struct {
	Port     string
	Password string
	JwtKey   string
}

type App struct {
	Cfg              *Config
	schedulerHandler *handler.SchedulerHandler
}

type Option func(*App)

func NewApp(opts ...Option) *App {
	s := &App{
		Cfg: &Config{
			Port:     os.Getenv("TODO_PORT"),
			Password: os.Getenv("TODO_PASSWORD"),
			JwtKey:   os.Getenv("JWT_KEY"),
		},
	}

	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithHandler(schedulerHandler *handler.SchedulerHandler) Option {
	return func(a *App) {
		a.schedulerHandler = schedulerHandler
	}
}
