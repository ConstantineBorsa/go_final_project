package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const tasksLimit = 50

func GetTasks(w http.ResponseWriter, r *http.Request) {

	log.Println("Entered getTasksHandler")
	// Задаем лимит на количество возвращаемых задач

	// Получаем текущую дату для фильтрации задач
	now := time.Now().Format("20060102")

	// Запрос к базе данных для получения задач
	tasks := []Task{}
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date >= ? ORDER BY date ASC LIMIT ?`
	err := DB.Select(&tasks, query, now, tasksLimit)
	if err != nil {
		erresponse := ErrorResponse{Error: "Ошибка запроса к базе данных"}
		sendErrorResponse(w, http.StatusInternalServerError, erresponse)

		log.Println("Error querying tasks:", err)
		return
	}

	// Формируем ответ
	response := TasksResponse{
		Tasks: tasks,
	}

	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Кодируем ответ в JSON и отправляем клиенту
	if err := json.NewEncoder(w).Encode(response); err != nil {
		erresponse := ErrorResponse{Error: "Ошибка кодирования ответа"}
		sendErrorResponse(w, http.StatusInternalServerError, erresponse)
		log.Println("Error encoding response:", err)
	}
}
