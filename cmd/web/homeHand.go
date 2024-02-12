package main

import (
	"as3/pkg/models"
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

func LatestArticles(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	articles, err := models.GetLatestArticles(db)
	if err != nil {
		log.Printf("Error getting latest articles: %v", err) // Добавьте логирование ошибки
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tpl, err := template.ParseFiles("ui/html/home.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err) // Добавьте логирование ошибки
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tpl.Execute(w, struct{ Articles []models.Article }{Articles: articles})
	if err != nil {
		log.Printf("Error executing template: %v", err) // Добавьте логирование ошибки
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
