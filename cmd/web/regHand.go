package main

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
)

func registerHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "GET":
		tpl := template.Must(template.ParseFiles("ui/html/register.html"))
		tpl.Execute(w, nil)
	case "POST":
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		fullName := r.FormValue("fullName")
		email := r.FormValue("email")
		password := r.FormValue("password")

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error while hashing password", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(`INSERT INTO users (full_name, email, hashed_password, role) VALUES (?, ?, ?, 'student')`, fullName, email, hashedPassword)
		if err != nil {
			http.Error(w, "Email already in use", http.StatusConflict)
			return
		}

		http.Redirect(w, r, "/login?message=registration_success", http.StatusSeeOther)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
