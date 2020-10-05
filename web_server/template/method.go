package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type User struct {
	ID     int
	Name   string
	Active bool
}

func (u *User) PrintActive(uppercase bool) string { // метод
	if !u.Active {
		return ""
	}

	output := "method says user " + u.Name + " active"

	if uppercase {
		return strings.ToUpper(output)
	}

	return output
}

func main() {
	tmpl, err := template.
		New("").
		ParseFiles("method.html")
	if err != nil {
		panic(err)
	}

	users := []User{
		User{1, "Vasily", true},
		User{2, "Ivan", false},
		User{3, "Dmitry", true},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "method.html",
			struct {
				Users []User
			}{
				users,
			})
		if err != nil {
			panic(err)
		}
	})

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
