package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func handleAddTask(w http.ResponseWriter, r *http.Request) {

	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// Проверяем обязательное поле Title
	if task.Title == "" {
		response := ErrorResponse{Error: "Title field is required"}
		sendErrorResponse(w, http.StatusBadRequest, response)
		return
	}

	// Проверяем формат даты и преобразуем в формат 20060102
	date := task.Date
	if date == "" {
		date = time.Now().Format("20060102")
	}

	parsedDate, err := time.Parse("20060102", date)
	if err != nil {
		response := ErrorResponse{Error: "Invalid date format, should be YYYYMMDD"}
		sendErrorResponse(w, http.StatusBadRequest, response)
		return
	}

	// Если дата задачи меньше текущей, вычисляем следующую дату выполнения
	if parsedDate.Before(time.Now()) {
		if task.Repeat != "" {
			nextDate, err := NextDate(time.Now(), date, task.Repeat)
			if err != nil {
				response := ErrorResponse{Error: fmt.Sprintf("Failed to calculate next date: %s", err.Error())}
				sendErrorResponse(w, http.StatusBadRequest, response)
				return
			}
			task.Date = nextDate
		} else {
			task.Date = time.Now().Format("20060102")
		}
	}

	// Выполняем SQL-запрос для добавления задачи
	insertSQL := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	result, err := DB.Exec(insertSQL, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		log.Printf("Failed to insert task into database: %v\n", err)
		response := ErrorResponse{Error: "Failed to insert task into database"}
		sendErrorResponse(w, http.StatusInternalServerError, response)
		return
	}

	// Получаем идентификатор добавленной задачи
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Failed to retrieve last insert ID: %v\n", err)
		response := ErrorResponse{Error: "Failed to retrieve last insert ID"}
		sendErrorResponse(w, http.StatusInternalServerError, response)
		return
	}

	// Возвращаем успешный ответ с ID
	response := map[string]int{"id": int(id)}
	sendResponse(w, http.StatusCreated, response)
}
