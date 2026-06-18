package services

import "github.com/Emin-07/final-project/internal/core/port"

type SchedulerServ struct {
	repo ports.SchedulerRepo
}

func NewSchedulerService(repo ports.SchedulerRepo) *SchedulerServ {
	return &SchedulerServ{repo: repo}
}
