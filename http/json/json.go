package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type UserMain struct {
	ID       int
	Username string
	phone    string
}

var jsonMainStr = `{"id": 42, "username": "eugene", "phone": "123"}`

func main() {
	data := []byte(jsonMainStr)

	u := &UserMain{}               // создаем структуру User
	err := json.Unmarshal(data, u) // распаковка data в парамер u
	if err != nil {                // обязательно проверять на ошибку
		// panic(err) - никогда  не надо использовать panic()
		// http.Error(w, "internal error", 500) - вывод ошибки для веб
		log.Fatalf("error")
	}
	fmt.Printf("struct:\n\t%#v\n\n", u)
	u.phone = "987654321"
	result, err := json.Marshal(u) // запаковка обратно структуры
	if err != nil {
		log.Fatalf("error")
	}
	//!!! endless
	fmt.Printf("json string:\n\t%s\n", string(result))
}
