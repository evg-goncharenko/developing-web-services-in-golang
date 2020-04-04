package main

import (
	"fmt"
	"net/http"
)

// action - куда слать запрос
var loginFormTmpl = []byte(`
<html>
	<body>
	<form action="/" method="post"> 
		Login: <input type="text" name="login">
		Password: <input type="password" name="password">
		<input type="submit" value="Login">
	</form>
	</body>
</html>
`)

func mainPagePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // если мы пришли не через пост-запрос
		w.Write(loginFormTmpl)
		return
	}

	// первый способ обработки параметров:
	// r.ParseForm() // вызов метода у реквеста
	// inputLogin := r.Form["login"][0] // распарсим тело запроса и добавим в Form, обращение к пост-параметрам

	// второй способ обработки параметров:
	inputLogin := r.FormValue("login") // FormValue() - просмотр всех параметров запроса
	inputPass := r.FormValue("password")
	fmt.Fprintln(w, "you enter: ", inputLogin, inputPass)
}

func main() {
	http.HandleFunc("/", mainPagePost)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
