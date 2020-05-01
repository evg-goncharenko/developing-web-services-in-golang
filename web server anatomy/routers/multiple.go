package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
)

func RegularRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Request with averange amout of logic\n")
}

func FastRequest(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprint(w, "Request with high hitrate\n")
}

func ComplexRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Request with complex routing logic\n")
}

func main() {
	// для быстрой обработки запросов
	fastApiHandler := httprouter.New()
	fastApiHandler.GET("/fast/:id", FastRequest)

	// для многофункциональной обработки запросов
	complexApiHandler := mux.NewRouter()
	complexApiHandler.HandleFunc("/complex/", ComplexRequest).
		Headers("X-Requested-With", "XMLHttpRequest") // ajax

	// использование стандартного роутера
	stdApiHandler := http.NewServeMux()
	stdApiHandler.HandleFunc("/std/", RegularRequest)

	siteMux := http.NewServeMux()
	siteMux.Handle("/fast/", fastApiHandler)
	siteMux.Handle("/complex/", complexApiHandler)
	siteMux.Handle("/std/", stdApiHandler)

	fmt.Println("starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", siteMux))
}
