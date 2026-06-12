package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

const maxDayAmount = 400

const weekDaysAmount = 7
const monthMaxDaysAmount = 31
const monthMinDaysAmount = 28
const monthsAmount = 12

func adjustedNegativeMonthDays(monthDays []int, date time.Time) []int {
	res := make([]int, 0, 31)
	for _, day := range monthDays {
		if day < 0 {
			days := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
			res = append(res, days+day)
		} else {
			res = append(res, day)
		}
	}
	return res
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	startTime, err := time.Parse(dateFormat, dstart)
	dateToAdd := map[string]int{
		"year": 0,
		"days": 0,
	}
	if err != nil {
		return "", err
	}
	if repeat == "" {
		return "", fmt.Errorf("repeat arg is empty")
	}
	repeatLetters := strings.Split(repeat, " ")
	if len(repeatLetters) < 2 && repeatLetters[0] != "y" {
		return "", fmt.Errorf("repeat format is incorrect")
	}
	switch repeatLetters[0] {
	case "d":
		daysToAdd, err := strconv.Atoi(repeatLetters[1])
		if err != nil {
			return "", err
		}
		if daysToAdd > maxDayAmount {
			return "", fmt.Errorf("you can only add up to %d amount of days, %d is too much")
		}
		dateToAdd["days"] = daysToAdd
	case "y":
		if len(repeatLetters) != 1 {
			return "", fmt.Errorf("too many arguments")
		}
		dateToAdd["year"] = 1

	/*
		TODO:
		There are some similarities between d,y and w,m. Maybe you should
		Find some common points and extract those into functions, to make
		code cleaner and better overall

		AND
		There are still bunch of error cases unhandled
	*/
	case "w":
		weekDaysString := strings.Split(repeatLetters[1], ",")
		weekDays := make([]int, 0, weekDaysAmount)
		days := map[string]int{
			"1": 1,
			"2": 2,
			"3": 3,
			"4": 4,
			"5": 5,
			"6": 6,
			"7": 0,
		}
		for _, weekDay := range weekDaysString {
			day, ok := days[weekDay]
			if !ok {
				return "", fmt.Errorf("day of the week supposed to be from 1 to %d, not %v", weekDaysAmount, weekDay)
			}
			weekDays = append(weekDays, day)
		}

		for {
			if startTime.After(now) && slices.Contains(weekDays, int(startTime.Weekday())) {
				return startTime.Format(dateFormat), nil
			}
			startTime = startTime.AddDate(0, 0, 1)
		}

	case "m":
		monthDaysString := strings.Split(repeatLetters[1], ",")
		monthDays := make([]int, 0, monthMaxDaysAmount)
		for _, monthDay := range monthDaysString {
			day, err := strconv.Atoi(monthDay)
			if err != nil {
				return "", err
			}
			if day > monthMaxDaysAmount || day == 0 || day < -monthMinDaysAmount {
				return "", fmt.Errorf("day of the month supposed to be from -%d to -1 and from 1 to %d. not %d", monthMinDaysAmount, monthMaxDaysAmount, day)
			}

			monthDays = append(monthDays, day)
		}
		months := make([]int, 0, monthsAmount)
		if len(repeatLetters) == 3 {
			monthsString := strings.Split(repeatLetters[2], ",")
			for _, month := range monthsString {
				monthN, err := strconv.Atoi(month)

				if err != nil {
					return "", err
				}
				if monthN > monthsAmount || monthN < 1 {
					return "", fmt.Errorf("month supposed to be from 1 to %d. not %d", monthsAmount, monthN)
				}
				months = append(months, monthN)
			}
		} else {
			months = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
		}

		for cnt := 1; ; cnt++ {
			if !slices.Contains(months, int(startTime.Month())) {
				days := time.Date(startTime.Year(), startTime.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
				startTime = startTime.AddDate(0, 0, days-startTime.Day()+1)
			} else {
				startTime = startTime.AddDate(0, 0, 1)
			}
			if startTime.After(now) && slices.Contains(adjustedNegativeMonthDays(monthDays, startTime), startTime.Day()) && slices.Contains(months, int(startTime.Month())) {
				return startTime.Format(dateFormat), nil
			}
		}
	default:
		return "", fmt.Errorf("got %s, expected one of [ d w m y ]", repeatLetters[0])
	}

	if now.Before(startTime) {
		startTime = startTime.AddDate(dateToAdd["year"], 0, dateToAdd["days"])
	}

	for now.After(startTime) {
		startTime = startTime.AddDate(dateToAdd["year"], 0, dateToAdd["days"])
	}

	return startTime.Format(dateFormat), nil
}
