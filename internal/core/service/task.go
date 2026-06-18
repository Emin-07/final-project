package service

import (
	"context"

	"github.com/Emin-07/final-project/internal/core/domain"
)

func (ss *SchedulerServ) AddTask(ctx context.Context, task *domain.Task) (int64, error) {
	id, err := ss.repo.Add(ctx, task)
	return id, err
}

func (ss *SchedulerServ) GetTask(ctx context.Context, id string) (*domain.Task, error) {
	return ss.repo.Get(ctx, id)
}

func (ss *SchedulerServ) ChangeTask(ctx context.Context, task *domain.Task) error {
	return ss.repo.Update(ctx, task)
}

func (ss *SchedulerServ) ChangeTaskDate(ctx context.Context, date string, id string) error {
	return ss.repo.UpdateDate(ctx, date, id)
}

func (ss *SchedulerServ) GetTasks(ctx context.Context, limit string, search string) ([]*domain.Task, error) {
	if limit == "" {
		limit = "50"
	}
	return ss.repo.Tasks(ctx, limit, search)

}

func (ss *SchedulerServ) CompleteTask(ctx context.Context, task *domain.Task) error {
	if task.Repeat == "" {
		err := ss.DeleteTask(ctx, task.ID)
		if err != nil {
			return err
		}
	} else {
		nextDate, err := ss.NextDate(GetNowYMD(), task.Date, task.Repeat)
		if err != nil {
			return err
		}
		err = ss.ChangeTaskDate(ctx, nextDate, task.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ss *SchedulerServ) DeleteTask(ctx context.Context, id string) error {
	return ss.repo.Delete(ctx, id)
}
