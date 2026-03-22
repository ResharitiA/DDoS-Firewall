package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Константа для имени куки (сессии)
const sessionCookieName = "session_token"

// Структура пользователя (ID с большой буквы для отображения в HTML)
type User struct {
	ID        int
	FirstName string
	LastName  string
	Phone     string
}

func init() {
	var err error
	// Открываем БД
	db, err = sql.Open("sqlite3", "./user_data.db")
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	// Создаем таблицы и дефолтного админа (admin/admin)
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        first_name TEXT,
        last_name TEXT,
        phone TEXT
    );
    CREATE TABLE IF NOT EXISTS admins (
        username TEXT PRIMARY KEY,
        password TEXT
    );
    INSERT OR IGNORE INTO admins (username, password) VALUES ('admin', 'admin');`

	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("Ошибка инициализации таблиц:", err)
	}
}

// Проверка: залогинен ли пользователь
func isAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie(sessionCookieName)
	return err == nil && cookie.Value == "authenticated"
}

func main() {
	// Маршруты
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/admin", adminHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/logout", logoutHandler)

	fmt.Println(">>> Сервер запущен: http://localhost:8080")
	fmt.Println(">>> Логин: admin | Пароль: admin")

	// Авто-открытие браузера через 1 секунду
	go func() {
		time.Sleep(1 * time.Second)
		openBrowser("http://localhost:8080/register")
	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// --- ОБРАБОТЧИКИ ---

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		phone := r.FormValue("phone")

		_, err := db.Exec("INSERT INTO users (first_name, last_name, phone) VALUES (?, ?, ?)", firstName, lastName, phone)
		if err != nil {
			http.Error(w, "Ошибка сохранения", 500)
			return
		}
		fmt.Fprintf(w, "Пользователь %s успешно зарегистрирован!", firstName)
	} else {
		tmpl := template.Must(template.ParseFiles("templates/register.html"))
		tmpl.Execute(w, nil)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var storedPass string
		err := db.QueryRow("SELECT password FROM admins WHERE username = ?", username).Scan(&storedPass)

		if err != nil || storedPass != password {
			http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
			return
		}

		// Ставим куку на 1 час
		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookieName,
			Value:    "authenticated",
			Path:     "/",
			HttpOnly: true,
		})
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	} else {
		tmpl := template.Must(template.ParseFiles("templates/login.html"))
		tmpl.Execute(w, nil)
	}
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	rows, err := db.Query("SELECT id, first_name, last_name, phone FROM users")
	if err != nil {
		http.Error(w, "Ошибка БД", 500)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Phone)
		users = append(users, u)
	}

	tmpl := template.Must(template.ParseFiles("templates/admin.html"))
	tmpl.Execute(w, users)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	id := r.URL.Query().Get("id")
	if id != "" {
		_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
		if err != nil {
			log.Println("Ошибка удаления:", err)
		}
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   sessionCookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Функция открытия браузера
func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	}
	if err != nil {
		log.Printf("Не удалось открыть браузер: %v", err)
	}
}