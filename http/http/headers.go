package main

import (
	"fmt"
	"net/http"
)

func handlerHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("RequestID", "d41d8cd98f00b204") // установка Header(), "параметр-значение"

	fmt.Fprintln(w, "You browser is", r.UserAgent())
	fmt.Fprintln(w, "You accept", r.Header.Get("Accept"))
	fmt.Fprintln(w, "You cookies", r.Header.Get("Cookie"))
}

func main() {
	http.HandleFunc("/", handlerHeaders)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
