package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID       int `json:"user_id,omitempty,string"`
	Username string
	Address  string `json:",omitempty"` // omitempty - если данных нет, то не выводить
	Company  string `json:"-"`          // поле Company не обрабатывается
}

func main() {
	u := &User{ // создаем структуру User
		ID:       1,
		Username: "eugene",
		Address:  "",
		Company:  "Mail.Ru Group",
	}
	result, _ := json.Marshal(u)
	fmt.Printf("json string: %s\n", string(result))
}
