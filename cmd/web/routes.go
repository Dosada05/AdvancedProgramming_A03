package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func requireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, authenticated := getUserRole(r)

			if !authenticated {
				http.Error(w, "Не аутентифицирован", http.StatusUnauthorized)
				return
			}

			roleAllowed := false
			for _, role := range allowedRoles {
				if userRole == role {
					roleAllowed = true
					break
				}
			}

			if !roleAllowed {
				http.Error(w, "Доступ запрещен", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getUserRole(r *http.Request) (string, bool) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Printf("Ошибка получения сессии: %v", err)
		return "", false
	}

	role, ok := session.Values["userRole"].(string)
	return role, ok
}

func router(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		LatestArticles(w, r, db)
	}).Methods("GET")

	// Используйте анонимные функции для удобного применения middleware
	createHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			ShowCreateArticleForm(w, r)
		} else if r.Method == "POST" {
			CreateArticle(w, r, db)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	// Применение middleware к маршруту создания статьи
	router.Handle("/create", requireRole("teacher", "admin")(createHandler))

	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		registerHandler(w, r, db)
	}).Methods("GET", "POST")

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		loginHandler(db, w, r)
	}).Methods("GET", "POST")

	// Админ-панель с применением middleware
	adminHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		adminPanelHandler(db, w, r)
	})
	router.Handle("/admin", requireRole("admin")(adminHandler)).Methods("GET")

	router.HandleFunc("/admin/change-role", func(w http.ResponseWriter, r *http.Request) {
		changeUserRoleHandler(db, w, r)
	}).Methods("POST")

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("ui/static/"))))

	return router
}
