package services

import (
	"fmt"
	"time"

	"github.com/Emin-07/final-project/internal/core/domain"
)

func (ss *SchedulerServ) ValidateStringTime(timeStr string) (time.Time, error) {
	if timeStr == "" {
		return ss.GetNowYMD(), nil
	}
	now, err := time.Parse(DateFormat, timeStr)
	if err != nil {
		return time.Date(0, 0, 0, 0, 0, 0, 0, time.Local), err
	}
	return now, nil
}

func (ss *SchedulerServ) ValidateDateFormat(date string) error {
	_, err := time.Parse(DateFormat, date)
	if err != nil {
		return fmt.Errorf("expected date format %v, got %v. err = %v", DateFormat, date, err)
	}
	return nil
}

func (ss *SchedulerServ) GetNowYMD() time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
}

func (ss *SchedulerServ) DateformatTime(t time.Time) string {
	return t.Format(DateFormat)
}

func (ss *SchedulerServ) CheckDate(task *domain.Task) error {
	now := ss.GetNowYMD()
	if task.Date == "" {
		task.Date = ss.DateformatTime(now)
	}
	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return err
	}

	if now.After(t) {
		if len(task.Repeat) == 0 {
			// если правила повторения нет, то берём сегодняшнее число
			task.Date = ss.DateformatTime(now)
		} else {
			// в противном случае, берём вычисленную ранее следующую дату
			next, err := ss.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
	}
	return nil
}
