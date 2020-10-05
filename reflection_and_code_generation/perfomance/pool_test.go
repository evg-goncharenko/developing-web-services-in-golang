package main

import (
	"bytes"
	"encoding/json"
	"sync"
	"testing"
)

const iterNum2 = 100

type PublicPage struct {
	ID          int
	Name        string
	Url         string
	OwnerID     int
	ImageUrl    string
	Tags        []string
	Description string
	Rules       string
}

// экземпляр структуры
var CoolGolangPublic = PublicPage{
	ID:          1,
	Name:        "CoolGolangPublic",
	Url:         "http://example.com",
	OwnerID:     100500,
	ImageUrl:    "http://example.com/img.png",
	Tags:        []string{"programming", "go", "golang"},
	Description: "Best page about golang programming",
	Rules:       "",
}

var Pages = []PublicPage{
	CoolGolangPublic,
	CoolGolangPublic,
	CoolGolangPublic,
}

// параллельная работа
func BenchmarkAllocNew(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			data := bytes.NewBuffer(make([]byte, 0, 512))
			_ = json.NewEncoder(data).Encode(Pages)
		}
	})
}

var dataPool = sync.Pool{ // один раз инициализируем
	New: func() interface{} { // определяем метод New
		return bytes.NewBuffer(make([]byte, 0, 512))
	},
}

func BenchmarkAllocPool(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			data := dataPool.Get().(*bytes.Buffer)
			_ = json.NewEncoder(data).Encode(Pages)
			data.Reset() // очистка данных во избежание утечки
			dataPool.Put(data)
		}
	})
}

/*
	Запуск:
	go test -bench . -benchmem pool_test.go
*/
