package main

import (
	"fmt"
	"net/http"
	"time"
)

// вывод различных префиксов:
func handlerMux(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "URL:", r.URL.String())
}

func adminPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ADMIN URL:", r.URL.String())
}

func adminPage2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ADMIN USER URL:", r.URL.String())
}

func main() {

	adminPrefix := "admin"

	admin := http.NewServeMux()           // создание нового роутера
	admin.HandleFunc("/user", adminPage2) // обработка юзера
	admin.HandleFunc("/", adminPage)      // обработка корня

	adminHandler := http.StripPrefix( // url: /admin/user --> /user
		"/"+adminPrefix,
		admin,
	)

	mux := http.NewServeMux()           // создание еще одного роутера
	mux.Handle("/admin/", adminHandler) // все что с префиксом /admin/ - обрабатывается через adminHandler
	mux.HandleFunc("/", handlerMux)     // все остальное - обрабатывается через handlerMux

	server := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second, // тайм-аут на чтение
		WriteTimeout: 10 * time.Second, // тайм-аут на запись
	}

	fmt.Println("starting server at :8080")
	server.ListenAndServe() // использование server вместо http
}
