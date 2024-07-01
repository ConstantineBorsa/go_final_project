package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем значение параметра id из запроса
	id := r.URL.Query().Get("id")

	// Проверяем, что id не пустой
	if id == "" {
		erresponse := ErrorResponse{Error: "Не указан идентификатор"}
		sendErrorResponse(w, http.StatusBadRequest, erresponse)
		return
	}

	// Запрос к базе данных для получения задачи по идентификатору
	var task Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := DB.Get(&task, query, id)
	if err != nil {
		erresponse := ErrorResponse{Error: "Задача не найдена"}
		sendErrorResponse(w, http.StatusNotFound, erresponse)
		log.Println("Error querying task:", err)
		return
	}

	// Формируем ответ
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		erresponse := ErrorResponse{Error: "Ошибка кодирования ответа"}
		sendErrorResponse(w, http.StatusInternalServerError, erresponse)
		log.Println("Error encoding response:", err)
		return
	}
}
