package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func taskAsDoneHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем значение параметра id из запроса
	id := r.URL.Query().Get("id")

	// Проверяем, что id не пустой
	if id == "" || id == "0" {
		erresponse := ErrorResponse{Error: "Не указан идентификатор"}
		sendErrorResponse(w, http.StatusBadRequest, erresponse)
		return
	}

	var err error

	// Получаем задачу из базы данных для дальнейших операций
	var task Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	if err := DB.Get(&task, query, id); err != nil {
		log.Printf("Failed to retrieve task from database: %v\n", err)
		response := ErrorResponse{Error: "Задача не найдена"}
		sendErrorResponse(w, http.StatusNotFound, response)
		return
	}

	// Если задача периодическая, вычисляем следующую дату выполнения
	if task.Repeat != "" {
		nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			response := ErrorResponse{Error: fmt.Sprintf("Failed to calculate next date: %s", err.Error())}
			sendErrorResponse(w, http.StatusInternalServerError, response)
			return
		}

		// Обновляем дату задачи
		updateSQL := `UPDATE scheduler SET date = ? WHERE id = ?`
		if _, err = DB.Exec(updateSQL, nextDate, id); err != nil {
			log.Printf("Failed to update task date in database: %v\n", err)
			response := ErrorResponse{Error: "Ошибка при обновлении даты задачи"}
			sendErrorResponse(w, http.StatusInternalServerError, response)
			return
		}
	} else {
		// Если значение repeat пустое, удаляем задачу из базы данных
		deleteSQL := `DELETE FROM scheduler WHERE id = ?`
		if _, err = DB.Exec(deleteSQL, id); err != nil {
			log.Printf("Failed to delete task from database: %v\n", err)
			response := ErrorResponse{Error: "Ошибка при удалении задачи"}
			sendErrorResponse(w, http.StatusInternalServerError, response)
			return
		}
	}

	// Успешный ответ
	sendResponse(w, http.StatusOK, map[string]interface{}{})
}
