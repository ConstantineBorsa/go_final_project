package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

func TaskAsDone(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем значение параметра id из запроса
	id := r.URL.Query().Get("id")

	// Проверяем, что id не пустой

	if id == "" {

		erresponse := ErrorResponse{Error: "Не указан идентификатор"}
		sendErrorResponse(w, http.StatusBadRequest, erresponse)
		return
	}

	var err error

	// Получаем задачу из базы данных для дальнейших операций
	var task Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	if err := db.Get(&task, query, id); err != nil {
		log.Printf("Failed to retrieve task from database: %v\n", err)
		response := ErrorResponse{Error: "Задача не найдена"}
		sendErrorResponse(w, http.StatusNotFound, response)
		return
	}

	now := time.Now()
	nowDate := now.Format("20060102")
	taskDate := task.Date

	// Если задача периодическая, вычисляем следующую дату выполнения
	if task.Repeat != "" {

		nextDate, err := nextDate(now, taskDate, task.Repeat)

		if err != nil {
			response := ErrorResponse{Error: fmt.Sprintf("Failed to calculate next date: %s", err.Error())}
			sendErrorResponse(w, http.StatusInternalServerError, response)
			return
		}

		// Если now и taskDate равны, добавляем интервал days к текущей дате
		if nowDate == taskDate {
			repeatParts := strings.Fields(task.Repeat)
			if len(repeatParts) == 2 {
				days, err := strconv.Atoi(repeatParts[1])
				if err == nil {
					nextDate = now.AddDate(0, 0, days).Format("20060102")
				}
			}
		}

		// Обновляем дату задачи
		updateSQL := `UPDATE scheduler SET date = ? WHERE id = ?`
		if _, err = db.Exec(updateSQL, nextDate, id); err != nil {
			log.Printf("Failed to update task date in database: %v\n", err)
			response := ErrorResponse{Error: "Ошибка при обновлении даты задачи"}
			sendErrorResponse(w, http.StatusInternalServerError, response)
			return
		}
	} else {
		// Если значение repeat пустое, удаляем задачу из базы данных
		deleteSQL := `DELETE FROM scheduler WHERE id = ?`
		if _, err = db.Exec(deleteSQL, id); err != nil {
			log.Printf("Failed to delete task from database: %v\n", err)
			response := ErrorResponse{Error: "Ошибка при удалении задачи"}
			sendErrorResponse(w, http.StatusInternalServerError, response)
			return
		}
	}

	// Успешный ответ
	sendResponse(w, http.StatusOK, map[string]interface{}{})
}
