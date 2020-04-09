// данный код относится к классу уязвимых программ
// пишем в терминал: curl -H 'User-Agent: <script>alert(1);<script>' http://127.0.0.1:8080
// опасность - Cross-Site Scripting
package main

import (
	"fmt"
	"net/http"
	"text/template" // библиотека для шаблонов
)

type tplParams struct {
	URL     string
	Browser string
}

// наполнение шаблонов данными (EXAMPLE или EXECUTE)
const EXAMPLE = `
Browser {{.Browser}}

you at {{.URL}}
`

var tmpl = template.New(`example`) // создание пустого шаблона

func handle(w http.ResponseWriter, r *http.Request) {
	params := tplParams{
		URL:     r.URL.String(),
		Browser: r.UserAgent(),
	}

	tmpl.Execute(w, params) // в переменную w запишется результат
}

func main() {
	tmpl, _ = tmpl.Parse(EXAMPLE) // лучше инициализировать шаблон в main (парсим его 1 раз)

	http.HandleFunc("/", handle)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
