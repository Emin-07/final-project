package handler

import (
	"github.com/Emin-07/final-project/internal/core/port"
)

type SchedulerHandler struct {
	service port.SchedulerService
}

func NewSchedulerHandler(service port.SchedulerService) *SchedulerHandler {
	return &SchedulerHandler{service: service}
}
