package main

import (
	"fmt"
	"net/http"
)

func runServer(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Addr:", addr, "URL:", r.URL.String())
		})

	// запуск отдельного веб-сервера
	mux.HandleFunc("/start", // отдельный url: /start
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Addr:", addr, "URL:", r.URL.String())
			addr := r.FormValue("addr") // получение параметров
			if addr != "" {
				go runServer(addr)
			}
		})

	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	fmt.Println("starting server at", addr)
	server.ListenAndServe() // запуск сервера - вызов метода структуры
}

func main() {
	go runServer(":8081")
	runServer(":8080") // запуск сервера в основной горутине
}

/*
	Запуск отдельного веб-сервера: `127.0.0.1:8081/start?addr=:8083`,
	а после можно: `127.0.0.1:8083`
*/
