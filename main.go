package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3" // Если используешь MySQL, замени на mysql драйвер
)

var db *sql.DB

func main() {
	var err error
	// Подключаемся к SQLite (файл из твоего дерева)
	db, err = sql.Open("sqlite3", "./cyberdefense_db.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создаем таблицу, если её нет
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name TEXT,
		last_name TEXT,
		phone TEXT
	)`)

	// 1. Эндпоинт для сохранения данных из игры
	http.HandleFunc("/submit_data", handleRegistration)

	// 2. Раздача статики (React билд)
	// Важно: Go будет искать файлы в папке retro-game/dist
	distPath := "./retro-game/dist"
	fs := http.FileServer(http.Dir(distPath))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, существует ли физический файл (css, js, png)
		path := distPath + r.URL.Path
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// Если файла нет (например, путь React-роутера), отдаем index.html
			http.ServeFile(w, r, distPath+"/index.html")
			return
		}
		fs.ServeHTTP(w, r)
	})

	fmt.Println("🚀 Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Читаем данные из FormData (как в твоем script.js)
	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	phone := r.FormValue("phone")

	// Пишем в SQL
	_, err := db.Exec("INSERT INTO users (first_name, last_name, phone) VALUES (?, ?, ?)", 
		firstName, lastName, phone)
	
	if err != nil {
		log.Println("Ошибка записи в БД:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf("✅ Зарегистрирован агент: %s %s (%s)\n", firstName, lastName, phone)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}