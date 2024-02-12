package main

import (
	"database/sql"
	"net/http"
)

func changeUserRoleHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	userID := r.FormValue("id")
	newRole := r.FormValue("role")

	_, err := db.Exec("UPDATE users SET role = ? WHERE id = ?", newRole, userID)
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
