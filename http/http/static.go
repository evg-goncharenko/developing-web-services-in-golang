package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
		<b> Hello World! </b> <br />
		<img src="/data/img/gopher.png" />
	`))
}

func main() {
	http.HandleFunc("/", handler)

	staticHandler := http.StripPrefix(
		"/data/",
		http.FileServer(http.Dir("./documents")), // FileServer теперь будет брать файлы из /Documents
	)
	http.Handle("/data/", staticHandler)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
