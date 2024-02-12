package main

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
)

func loginHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tpl := template.Must(template.ParseFiles("ui/html/login.html"))
		tpl.Execute(w, nil)
	case "POST":
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		// Проверить учетные данные пользователя
		var (
			userID                   int
			hashedPassword, userRole string
		)
		err = db.QueryRow("SELECT id, hashed_password, role FROM users WHERE email = ?", email).Scan(&userID, &hashedPassword, &userRole)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		// Сравнить предоставленный пароль с хешированным паролем
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		session, _ := store.Get(r, "session-name")
		session.Values["userID"] = userID
		session.Values["userRole"] = userRole
		session.Save(r, w)

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
