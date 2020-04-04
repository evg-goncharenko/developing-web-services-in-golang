package main

import (
	"fmt"
	"net/http"
)

func handlerGet(w http.ResponseWriter, r *http.Request) {
	myParam := r.URL.Query().Get("param") // первый способ получить параметры (вернется строка)
	if myParam != "" {
		fmt.Fprintln(w, "`myParam` is", myParam)
	}

	key := r.FormValue("key") // второй способ получить параметры
	if key != "" {
		fmt.Fprintln(w, "`key` is", key)
	}
}

func main() {
	http.HandleFunc("/", handlerGet)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

/*
	127.0.0.1:8080/?param=123&key=test
*/
