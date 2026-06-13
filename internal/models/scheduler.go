package models

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"time"
)

const selectAllFromScheduler = `SELECT id, date, title, comment, repeat FROM scheduler `

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

func (s *SchedulerModel) Tasks(ctx context.Context, limit int, search string) ([]*Task, error) {
	var err error
	var rows *sql.Rows
	if search == "" {
		query := selectAllFromScheduler + `ORDER BY date LIMIT :limit`
		rows, err = s.DB.QueryContext(ctx, query, sql.Named("limit", limit))
	} else {
		if timeSearch, err := time.Parse("02.01.2006", search); err == nil {
			date := timeSearch.Format("20060102")
			query := selectAllFromScheduler + `WHERE date = :date LIMIT :limit`
			rows, err = s.DB.QueryContext(ctx, query, sql.Named("date", date), sql.Named("limit", limit))
		} else {
			query := selectAllFromScheduler + `WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit`
			rows, err = s.DB.QueryContext(ctx, query, sql.Named("search", fmt.Sprintf("%%%s%%", search)), sql.Named("limit", limit))
		}
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task

	for rows.Next() {
		task := &Task{}
		if err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	if tasks == nil {
		return []*Task{}, nil
	}
	return tasks, nil
}
