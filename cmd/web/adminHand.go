package main

import (
	"database/sql"
	"html/template"
	"net/http"
)

func adminPanelHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query("SELECT id, full_name, email, role FROM users")
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []struct {
		ID       int
		FullName string
		Email    string
		Role     string
	}

	for rows.Next() {
		var u struct {
			ID       int
			FullName string
			Email    string
			Role     string
		}
		if err := rows.Scan(&u.ID, &u.FullName, &u.Email, &u.Role); err != nil {
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	tpl := template.Must(template.ParseFiles("ui/html/admin_panel.html"))
	tpl.Execute(w, map[string]interface{}{"Users": users})
}
