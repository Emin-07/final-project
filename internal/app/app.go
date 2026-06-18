package app

import (
	"os"

	"github.com/Emin-07/final-project/internal/adapter/handler"
	"github.com/Emin-07/final-project/internal/core/port"
)

type Config struct {
	Port     string
	Password string
	JwtKey   string
}

type App struct {
	Cfg                  *Config
	schedulerHandler     *handler.SchedulerHandler
	schedulerServicePort ports.SchedulerService
	schedulerRepoPort    ports.SchedulerRepo
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

func WithService(schedulerServicePort ports.SchedulerService) Option {
	return func(a *App) {
		a.schedulerServicePort = schedulerServicePort
	}
}
func WithRepo(schedulerRepoPort ports.SchedulerRepo) Option {
	return func(a *App) {
		a.schedulerRepoPort = schedulerRepoPort
	}
}
