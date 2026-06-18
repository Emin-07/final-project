package service

import (
	"time"

	"github.com/Emin-07/final-project/internal/core/domain"
)

func ValidateStringTime(timeStr string) (time.Time, error) {
	if timeStr == "" {
		return GetNowYMD(), nil
	}
	now, err := time.Parse(DateFormat, timeStr)
	if err != nil {
		return time.Date(0, 0, 0, 0, 0, 0, 0, time.Local), err
	}
	return now, nil
}

func GetNowYMD() time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
}

func DateformatTime(t time.Time) string {
	return t.Format(DateFormat)
}

func (ss *SchedulerServ) CheckDate(task *domain.Task) error {
	now := GetNowYMD()
	if task.Date == "" {
		task.Date = DateformatTime(now)
	}
	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return err
	}

	if now.After(t) {
		if len(task.Repeat) == 0 {
			// если правила повторения нет, то берём сегодняшнее число
			task.Date = DateformatTime(now)
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
