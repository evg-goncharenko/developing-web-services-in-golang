package main

// не используем какие-нибудь внешние пакеты
import (
	"api/handlers"
	"net/http"
	"sync"
)

// GET - получение
// POST - добавление новых данных (при повторении запроса должно измениться состояние сервера)
// PUT - изменение данных (идемпотентен)
// DELETE - удаление

// HEAD - запрашивает ресурс так же, как и метод GET, но без тела ответа
// PATCH - запрос на изменение данных, но без записи в хранилище (возможна передача одного поля)
// OPTIONS - изменение настроек сервера

func main() {

	users := map[string]*handlers.User{
		"test": &handlers.User{
			ID:       1,
			Login:    "test",
			Password: "test",
		},
	}

	sessions := map[string]*handlers.User{
		"tokenknsjkdfklsdf": users["test"],
	}

	mu := &sync.Mutex{}

	handler := handlers.Handler{
		Sessions: sessions,
		Users:    users,
		Mu:       mu,
	}

	http.HandleFunc("/users/", handler.HandleUsers)
	http.HandleFunc("/session/", handler.HandleSession)

	http.ListenAndServe(":8080", nil)
}
