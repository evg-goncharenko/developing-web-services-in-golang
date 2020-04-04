package main

import (
	"fmt"
	"net/http"
)

func handlerPages(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Main page")
}

// HandleFunc() - роутер
func main() {
	// url выберется самый длинный среди всех подходящих
	http.HandleFunc("/page",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Single page:", r.URL.String())
		})

	http.HandleFunc("/collection/test",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Single page collection:", r.URL.String())
		})

	http.HandleFunc("/collection/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Multiple pages:", r.URL.String())
		})

	http.HandleFunc("/", handlerPages) // в него пойдет все, что не попало выше

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
