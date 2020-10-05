package main

import (
	"fmt"
	"html/template" // экранирование users
	"net/http"
)

type User struct {
	ID     int
	Name   string
	Active bool
}

func main() {
	// Must() принимает распарсенный файл, если их нет - будет паника
	tmpl := template.Must(template.ParseFiles("users.html"))

	users := []User{
		User{1, "Eugene", true},
		User{2, "<script>alert(document.cookie)</script>", false},
		User{3, "Edgar", true},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w,
			struct {
				Users []User
			}{
				users,
			})
	})

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
