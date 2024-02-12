package main

import (
	"as3/pkg/mysql"
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("key"))

func main() {
	db, err := mysql.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := router(db)

	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
