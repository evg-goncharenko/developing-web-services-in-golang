package main

import (
	"encoding/json"
	"fmt"
)

var jsonStr = `[
	{"id": 17, "username": "eugene4", "phone": 0},
	{"id": "17", "address": "none", "company": "Mail.ru"}
]`

func main() {
	data := []byte(jsonStr)

	var user1 interface{}
	json.Unmarshal(data, &user1) // type(user1) = []interface {} {map[string] interface{}, map[string] interface{}}
	fmt.Printf("unpacked in empty interface:\n%#v\n\n", user1)

	user2 := map[string]interface{}{
		"id":       42,
		"username": "eugene",
	}
	// var user2i interface{} = user2
	result, _ := json.Marshal(user2)
	fmt.Printf("json string from map:\n %s\n", string(result))
}
