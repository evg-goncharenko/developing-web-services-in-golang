package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

// порядок полей в структуре важен!
type User struct {
	Login    string
	RealName string `unpack:"-" json:"real_name"` // теги: unpack и json
	ID       int
	Flags    int
}

type Unpacker interface {
	Unpack([]byte) error // метод Unpack()
}

func UnpackReflect(u interface{}, data []byte) error {

	if unp, ok := u.(Unpacker); ok {
		return unp.Unpack(data)
	}

	r := bytes.NewReader(data) // стандартный метод

	val := reflect.ValueOf(u).Elem() // получение значения

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		if typeField.Tag.Get("unpack") == "-" {
			continue
		}

		switch typeField.Type.Kind() { // проверка на тип поля
		case reflect.Int:
			var value uint32
			binary.Read(r, binary.LittleEndian, &value)
			valueField.Set(reflect.ValueOf(int(value)))
		case reflect.String:
			var lenRaw uint32
			binary.Read(r, binary.LittleEndian, &lenRaw)

			dataRaw := make([]byte, lenRaw)
			binary.Read(r, binary.LittleEndian, &dataRaw)

			valueField.SetString(string(dataRaw))
		default:
			return fmt.Errorf("bad type: %v for field %v", typeField.Type.Kind(), typeField.Name)
		}
	}

	return nil
}

// Intel процессоры всегда работают в формате .LittleEndian (младшие разряды в начале)
func main() {
	/*
		perl -E '$b = pack("L L/a* L", 1_123_456, "d.dorofeev", 16);
			print map { ord.", "  } split("", $b); '
	*/
	data := []byte{ // бинарное хранение данных

		10, 0, 0, 0, // 4 байта - int, 10 - кол-во символов в строке
		100, 46, 100, 111, 114, 111, 102, 101, 101, 118, // строка

		128, 36, 17, 0,

		16, 0, 0, 0,
	}
	u := new(User)
	err := UnpackReflect(u, data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", u)
}
