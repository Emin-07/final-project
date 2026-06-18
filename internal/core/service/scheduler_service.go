package service

import "github.com/Emin-07/final-project/internal/core/port"

type SchedulerServ struct {
	repo port.SchedulerRepo
}

func NewSchedulerService(repo port.SchedulerRepo) *SchedulerServ {
	return &SchedulerServ{repo: repo}
}
