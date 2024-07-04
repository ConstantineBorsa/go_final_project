package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {

	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// Проверяем обязательное поле ID
	if task.ID == "" {
		response := ErrorResponse{Error: "Не указан идентификатор задачи (ID)"}
		sendErrorResponse(w, http.StatusBadRequest, response)
		return
	}

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

	// Выполняем SQL-запрос для обновления задачи
	updateSQL := `UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`
	result, err := DB.Exec(updateSQL, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		log.Printf("Failed to update task in database: %v\n", err)
		response := ErrorResponse{Error: "Failed to update task in database"}
		sendErrorResponse(w, http.StatusInternalServerError, response)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to get rows affected: %v\n", err)
		response := ErrorResponse{Error: "Failed to get rows affected"}
		sendErrorResponse(w, http.StatusInternalServerError, response)
		return
	}

	// Проверяем, была ли обновлена хотя бы одна запись
	if rowsAffected == 0 {
		response := ErrorResponse{Error: "Задача не найдена"}
		sendErrorResponse(w, http.StatusNotFound, response)
		return
	}

	// Возвращаем успешный ответ
	sendResponse(w, http.StatusOK, map[string]string{})
}
