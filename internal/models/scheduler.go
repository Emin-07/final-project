package models

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
)

//go:embed create_schedule_table.sql
var createTableSQL string

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// TODO: Make interface to make it abstract and easier to test 'just to flex'

type SchedulerModel struct {
	DB *sql.DB
}

func (s *SchedulerModel) InitScheduleTable(ctx context.Context) error {
	_, err := s.DB.ExecContext(ctx, createTableSQL)
	return err
}

func (s *SchedulerModel) AddTask(ctx context.Context, task *Task) (int64, error) {
	query := `INSERT INTO scheduler(date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
	res, err := s.DB.ExecContext(ctx, query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, fmt.Errorf("insert task: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("last insert id: %w", err)
	}
	return id, nil
}
