package main

import (
	"fmt"
	"net/http"
	"time"
)

func mainPageCookies(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id") // вызываем метод Cookie()
	loggedIn := err != http.ErrNoCookie    // проверка если Cookie() нет

	if loggedIn { // если мы залогинены
		fmt.Fprintln(w, `<a href="/logout">logout</a>`)
		fmt.Fprintln(w, "Welcome, "+session.Value) // выводим значение из Cookie()
	} else {
		fmt.Fprintln(w, `<a href="/login">login</a>`)
		fmt.Fprintln(w, "You need to login")
	}
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().Add(10 * time.Hour) // время через которое Cookie() завершится (через 10 часов)
	cookie := http.Cookie{                       // создаем новую структуру Cookie()
		Name:    "session_id",
		Value:   "eugene",
		Expires: expiration, // время действия
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound) // перенаправление пользователя на главную страницу
}

func logoutPage(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1) // удаление Cookie()
	http.SetCookie(w, session)

	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/logout", logoutPage)
	http.HandleFunc("/", mainPageCookies)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
