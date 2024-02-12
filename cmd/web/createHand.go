package main

import (
	"as3/pkg/models"
	"database/sql"
	"html/template"
	"net/http"
)

func ShowCreateArticleForm(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("ui/html/create_article.html"))
	tpl.Execute(w, nil)
}

func CreateArticle(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	article := models.Article{
		Title:   title,
		Content: content,
	}

	err := models.CreateArticle(db, article)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
