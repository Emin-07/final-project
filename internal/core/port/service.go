package port

import (
	"context"
	"time"

	"github.com/Emin-07/final-project/internal/core/domain"
)

type SchedulerService interface {
	NextDate(now time.Time, dstart string, repeat string) (string, error)
	CheckDate(task *domain.Task) error

	AddTask(ctx context.Context, task *domain.Task) (int64, error)
	GetTask(ctx context.Context, id string) (*domain.Task, error)
	ChangeTask(ctx context.Context, task *domain.Task) error
	ChangeTaskDate(ctx context.Context, date string, id string) error
	GetTasks(ctx context.Context, limit string, search string) ([]*domain.Task, error)
	CompleteTask(ctx context.Context, task *domain.Task) error
	DeleteTask(ctx context.Context, id string) error
}
