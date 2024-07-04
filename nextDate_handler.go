package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	if nowStr == "" || dateStr == "" || repeat == "" {
		http.Error(w, "missing query parameters", http.StatusBadRequest)
		return
	}

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "invalid now date format", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, nextDate)
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	// Парсим исходную дату
	startDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("invalid date format")
	}

	// Проверяем правило повторения
	if repeat == "" {
		return "", errors.New("repeat rule is empty")
	}

	repeatParts := strings.Fields(repeat)
	if len(repeatParts) < 1 {
		return "", errors.New("invalid repeat rule")
	}

	repeatType := repeatParts[0]

	// Определяем ближайшую дату на основе правила повторения
	var days int
	var nextDate time.Time
	switch repeatType {
	case "y":
		if len(repeatParts) != 1 {
			return "", errors.New("invalid year repeat rule")
		}
		nextDate = startDate.AddDate(1, 0, 0)
	case "d":
		if len(repeatParts) != 2 {
			return "", errors.New("invalid day repeat rule")
		}

		days, err = strconv.Atoi(repeatParts[1])

		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("invalid day interval")
		}
		nextDate = startDate.AddDate(0, 0, days)

	default:
		return "", errors.New("unsupported repeat rule")
	}

	// Если дата равна сегодняшней, возвращаем сегодняшнюю дату
	if time.Now().Format("20060102") == startDate.Format("20060102") {
		return time.Now().Format("20060102"), nil
	}

	// Если вычисленная дата не больше текущей, увеличиваем её на интервал повторения
	for !nextDate.After(now) {
		switch repeatType {
		case "y":
			nextDate = nextDate.AddDate(1, 0, 0)
		case "d":
			nextDate = nextDate.AddDate(0, 0, days)

		}
	}

	return nextDate.Format("20060102"), nil
}
