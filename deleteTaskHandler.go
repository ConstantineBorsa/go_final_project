package main

import (
	"log"
	"net/http"
)

func deleteTask(w http.ResponseWriter, r *http.Request) {
	// Получаем значение параметра id из запроса
	id := r.URL.Query().Get("id")

	// Проверяем, что id не пустой
	if id == "" || id == "0" {
		erresponse := ErrorResponse{Error: "Не указан идентификатор"}
		sendErrorResponse(w, http.StatusBadRequest, erresponse)
		return
	}

	// Удаляем задачу из базы данных
	deleteSQL := `DELETE FROM scheduler WHERE id = ?`
	if _, err := DB.Exec(deleteSQL, id); err != nil {
		log.Printf("Failed to delete task from database: %v\n", err)
		response := ErrorResponse{Error: "Ошибка при удалении задачи"}
		sendErrorResponse(w, http.StatusInternalServerError, response)
		return
	}

	// Успешный ответ
	sendResponse(w, http.StatusOK, map[string]interface{}{})
}
