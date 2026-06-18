package handler

import (
	"github.com/Emin-07/final-project/internal/core/port"
)

type SchedulerHandler struct {
	service ports.SchedulerService
}

func NewSchedulerHandler(service ports.SchedulerService) *SchedulerHandler {
	return &SchedulerHandler{service: service}
}
