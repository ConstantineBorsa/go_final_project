package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	h "github.com/ConstantineBorsa/go_final_project/handlers"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const defaultPort = "7540"

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

	db, err := sqlx.Connect("sqlite3", dbFile)
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

		_, err = db.Exec(createTableSQL)
		if err != nil {
			log.Fatal(err)
		}

		createIndexSQL := `CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);`

		_, err = db.Exec(createIndexSQL)
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

	http.HandleFunc("/api/nextdate", h.NextDate)
	http.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		h.GetTasks(w, r, db)
	})
	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetTask(w, r, db)
		case http.MethodPost:
			h.AddTask(w, r, db)
		case http.MethodPut:
			h.UpdateTask(w, r, db)
		case http.MethodDelete:

			h.DeleteTask(w, r, db)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/api/task/done", func(w http.ResponseWriter, r *http.Request) {
		h.TaskAsDone(w, r, db)
	})

	http.Handle("/", fs)

	addr := os.Getenv("TODO_PORT")
	if addr == "" {
		addr = defaultPort
	}
	addr = ":" + addr

	log.Printf("Запуск сервера на порту %s...\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v\n", err)
	}
}
