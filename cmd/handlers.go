package main

import (
	"net/http"
	"time"
)

func (app *application) getNextDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	if date == "" {
		http.Error(w, "no date passed in parameters", http.StatusBadRequest)
		return
	}
	if repeat == "" {
		http.Error(w, "no repeat passed in parameters", http.StatusBadRequest)
		return
	}
	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now().UTC()
	} else {
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	res, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}
	w.Write([]byte(res))
}
