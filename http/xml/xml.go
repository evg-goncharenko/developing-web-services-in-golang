package main

import (
	"encoding/xml"
	"fmt"
)

type User struct {
	ID      int    `xml:"id,attr"` // чтение id не из атрибута (не из <id></id>)
	Login   string `xml:"login"`
	Name    string `xml:"name"`
	Browser string `xml:"browser"`
}

type Users struct { // верхнеуровневая структура
	Version string `xml:"version,attr"`
	List    []User `xml:"user"`
}

var xmlData = []byte(`<?xml version="1.0" encoding="utf-8"?>
<users>
	<user id="1">
		<login>user1</login>
		<name>Евгений Гончаренко</name>
		<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36
</browser>
	</user>
	<user id="2">
		<login>user2</login>
		<name>Иван Иванов</name>
		<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36
</browser>
	</user>
	<user id="2">
		<login>user3</login>
		<name>Иван Петров</name>
		<browser>Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0; Trident/5.0)</browser>
	</user>
	<user id="1">
		<login>user1</login>
		<name>Евгений Гончаренко</name>
		<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36
</browser>
	</user>
	<user id="2">
		<login>user2</login>
		<name>Иван Иванов</name>
		<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36
</browser>
	</user>
	<user id="2">
		<login>user3</login>
		<name>Иван Петров</name>
		<browser>Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0; Trident/5.0)</browser>
	</user>
</users>`)

func main() {
	users := new(Users)
	err := xml.Unmarshal(xmlData, &users)
	if err != nil {
		panic(err)
	}
	fmt.Println(users)

	res, _ := xml.Marshal(users.List[0])
	fmt.Println(string(res))
}
