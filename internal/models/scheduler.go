package models

import (
	"context"
	"database/sql"
	_ "embed"
)

//go:embed create_schedule_table.sql
var createTableSQL string

type Scheduler struct {
}

// TODO: Make interface to make it abstract and easier to test 'just to flex'
type SchedulerModel struct {
	DB *sql.DB
}

func (s *SchedulerModel) InitScheduleTable(ctx context.Context) error {
	_, err := s.DB.ExecContext(ctx, createTableSQL)
	return err
}
