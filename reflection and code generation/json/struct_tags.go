package main

import (
	"encoding/json"
	"fmt"
)

// пример реализации тегов структур:
type User struct {
	ID       int `json:"user_id,string"` // формат: `ключ:значение`
	Username string
	Address  string `json:",omitempty"`
	Company  string `json:"-"`
}

func main() {
	u := &User{
		ID:       42,
		Username: "rvasily",
		Address:  "test",
		Company:  "Mail.Ru Group",
	}
	result, _ := json.Marshal(u)
	fmt.Printf("json string: %s\n", string(result))
}
