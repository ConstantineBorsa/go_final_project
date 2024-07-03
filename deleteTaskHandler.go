package main

import (
	"log"
	"net/http"
	"strconv"
)

func deleteTask(w http.ResponseWriter, r *http.Request) {
	// Получаем значение параметра id из запроса
	idStr := r.URL.Query().Get("id")

	// Проверяем, что id не пустой
	if idStr == "" || idStr == "0" {
		erresponse := ErrorResponse{Error: "Не указан идентификатор"}
		sendErrorResponse(w, http.StatusBadRequest, erresponse)
		return
	}

	// Преобразуем id в число
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		erresponse := ErrorResponse{Error: "Некорректный идентификатор"}
		sendErrorResponse(w, http.StatusBadRequest, erresponse)
		return
	}

	// Удаляем задачу из базы данных
	deleteSQL := `DELETE FROM scheduler WHERE id = ?`
	result, err := DB.Exec(deleteSQL, id)
	if err != nil {
		log.Printf("Failed to delete task from database: %v\n", err)
		response := ErrorResponse{Error: "Ошибка при удалении задачи"}
		sendErrorResponse(w, http.StatusInternalServerError, response)
		return
	}

	// Проверяем количество удаленных строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to get rows affected: %v\n", err)
		response := ErrorResponse{Error: "Ошибка при удалении задачи"}
		sendErrorResponse(w, http.StatusInternalServerError, response)
		return
	}

	if rowsAffected == 0 {
		response := ErrorResponse{Error: "Задача не найдена"}
		sendErrorResponse(w, http.StatusNotFound, response)
		return
	}

	// Успешный ответ
	sendResponse(w, http.StatusOK, map[string]interface{}{})
}
