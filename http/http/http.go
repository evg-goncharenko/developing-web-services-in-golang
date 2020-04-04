package main

import (
	"fmt"
	"net/http"
)

// ResponseWriter - куда пишем результат, *http.Request - наш запрос
func handlerHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Привет, мир!") // 1-й способ вывода
	w.Write([]byte("!!!"))          // 2-й способ вывода
}

func main() {
	http.HandleFunc("/", handlerHTTP) // обработка запросов с вызовом необходимой функции (роутер)

	fmt.Println("starting server at :8080")
	// :8080 - адрес на котором слушаем, nil - используем роутер, который мы передали через http
	http.ListenAndServe(":8080", nil) // запуск дефолтного http сервера
}
