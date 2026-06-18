package scheduler

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/Emin-07/final-project/internal/core/domain"
)

const (
	selectAllFromScheduler = "SELECT id, date, title, comment, repeat FROM scheduler "
	whereId                = " WHERE id = :id "
	limitQuery             = " LIMIT :limit "
)

//go:embed create_schedule_table.sql
var createTableSQL string

func (s *SchedulerRepo) Add(ctx context.Context, task *domain.Task) (int64, error) {
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

func (s *SchedulerRepo) Tasks(ctx context.Context, limit string, search string) ([]*domain.Task, error) {
	var err error
	var rows *sql.Rows
	if search == "" {
		query := selectAllFromScheduler + `ORDER BY date` + limitQuery
		rows, err = s.DB.QueryContext(ctx, query, sql.Named("limit", limit))
	} else {
		if timeSearch, err := time.Parse("02.01.2006", search); err == nil {
			date := timeSearch.Format("20060102")
			query := selectAllFromScheduler + `WHERE date = :date` + limitQuery
			rows, err = s.DB.QueryContext(ctx, query, sql.Named("date", date), sql.Named("limit", limit))
		} else {
			query := selectAllFromScheduler + `WHERE title LIKE :search OR comment LIKE :search ORDER BY date` + limitQuery
			rows, err = s.DB.QueryContext(ctx, query, sql.Named("search", fmt.Sprintf("%%%s%%", search)), sql.Named("limit", limit))
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*domain.Task

	for rows.Next() {
		task := &domain.Task{}
		if err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	if tasks == nil {
		return []*domain.Task{}, nil
	}
	return tasks, nil
}

func (s *SchedulerRepo) Get(ctx context.Context, id string) (*domain.Task, error) {
	query := selectAllFromScheduler + whereId
	task := &domain.Task{}
	err := s.DB.QueryRowContext(ctx, query, sql.Named("id", id)).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNoRecord
		}
		return nil, err
	}
	return task, nil
}

func (s *SchedulerRepo) Update(ctx context.Context, task *domain.Task) error {
	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat` + whereId
	res, err := s.DB.ExecContext(ctx, query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID),
	)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}
	return nil
}

func (s *SchedulerRepo) UpdateDate(ctx context.Context, next string, id string) error {
	query := `UPDATE scheduler SET date = :date ` + whereId
	res, err := s.DB.ExecContext(ctx, query, sql.Named("date", next), sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("update date task: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}
	return nil
}

func (s *SchedulerRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM scheduler ` + whereId
	res, err := s.DB.ExecContext(ctx, query, sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}
	return nil

}

func HasNoTables(ctx context.Context, db *sql.DB, query string) (bool, error) {
	var count int
	err := db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func InitScheduleTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, createTableSQL)
	return err
}

func (s *SchedulerRepo) InitDataIntoDb(db *sql.DB) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	isEmpty, err := HasNoTables(ctx, db, `SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%'`)
	if err != nil {
		return err
	}

	if isEmpty {
		err = InitScheduleTable(ctx, db)
		if err != nil {
			return err
		}
	}

	return nil
}
