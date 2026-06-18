package ports

import (
	"context"

	"github.com/Emin-07/final-project/internal/core/domain"
)

type SchedulerRepo interface {
	Add(ctx context.Context, task *domain.Task) (int64, error)
	Tasks(ctx context.Context, limit string, search string) ([]*domain.Task, error)
	Get(ctx context.Context, id string) (*domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
	UpdateDate(ctx context.Context, next string, id string) error
	Delete(ctx context.Context, id string) error
}
