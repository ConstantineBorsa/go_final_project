package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sqlx.DB

func main() {

	// Определение пути к базе данных
	var dbFile string

	// Проверяем переменную окружения TODO_DBFILE
	envDBFile := os.Getenv("TODO_DBFILE")
	if envDBFile != "" {
		dbFile = envDBFile
	} else {
		// Если переменная не установлена, используем стандартный путь
		appPath, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		dbFile = filepath.Join(filepath.Dir(appPath), "scheduler.db")
	}

	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
		log.Println(dbFile)
		log.Println(install)
	}
	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX

	DB, err = sqlx.Connect("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	//defer DB.Close()

	if install {
		createTableSQL := `CREATE TABLE IF NOT EXISTS scheduler (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date TEXT NOT NULL,
            title TEXT NOT NULL,
            comment TEXT,
            repeat TEXT CHECK (length(repeat) <= 128)
        );`

		_, err = DB.Exec(createTableSQL)
		if err != nil {
			log.Fatal(err)
		}

		createIndexSQL := `CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);`

		_, err = DB.Exec(createIndexSQL)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Database and table created successfully")
	} else {
		log.Println("Database already exists")
	}

	// Устанавливаем директорию для файлов
	webDir := "./web"
	fs := http.FileServer(http.Dir(webDir))

	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/tasks", getTasksHandler)
	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTaskHandler(w, r)
		case http.MethodPost:
			handleAddTask(w, r)
		case http.MethodPut:
			updateTaskHandler(w, r)
		case http.MethodDelete:

			deleteTask(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/api/task/done", taskAsDoneHandler)

	http.Handle("/", fs)

	port1 := os.Getenv("TODO_PORT")
	if port1 == "" {
		port1 = "7540"
	}
	port1 = ":" + port1

	log.Printf("Запуск сервера на порту %s...\n", port1)
	if err := http.ListenAndServe(port1, nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v\n", err)
	}
}
