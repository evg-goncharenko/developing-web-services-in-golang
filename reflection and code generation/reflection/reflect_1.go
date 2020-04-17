package main

import (
	"fmt"
	"reflect"
)

type Login string

// благодаря reflection работает json-пакет
type User struct {
	ID       int
	RealName string `unpack:"-"`
	Login    Login
	Flags    int
}

func PrintReflect(u interface{}) error {
	// reflect. - метатип
	val := reflect.ValueOf(u) // получение метаинформации о типе

	// необходимо для работы с указателями на структуру
	if val.Kind() == reflect.Ptr { // является ли этот тип указателем
		val = val.Elem() // получаем значение
	}

	fmt.Println(val.Type(), val.Kind())

	// NumField() - кол-во полей
	fmt.Printf("%T have %d fields:\n", u, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)       // значение поля
		typeField := val.Type().Field(i) // тип поля

		fmt.Printf("\tname=%v, type=%v, value=%v, tag=`%v`\n",
			typeField.Name,
			typeField.Type,
			valueField,
			typeField.Tag,
		)
	}
	return nil
}

func main() {
	u := User{
		ID:       42,
		RealName: "evg",
		Flags:    32,
	}
	fmt.Printf("%#v", u)
	err := PrintReflect(u)
	if err != nil {
		panic(err)
	}
}
