package scheduler

import "database/sql"

type SchedulerRepo struct {
	DB *sql.DB
}

func NewSchedulerRepo(db *sql.DB) *SchedulerRepo {
	return &SchedulerRepo{DB: db}
}
