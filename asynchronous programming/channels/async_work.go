package main

import (
	"fmt"
	"time"
)

// загрузка нескольких типов данных
func getComments() <-chan string { // <-chan означает, что в этот канал мы будем только писать
	result := make(chan string, 1) // надо использовать буферизированный канал, тк не знаем, когда начнется чтение
	go func(out chan<- string) {   // асинхронный запрос
		time.Sleep(2 * time.Second)
		fmt.Println("async operation ready, return comments")
		out <- "32 comments"
	}(result) // возвращает канал
	return result
}

func getUser() <-chan string {
	result := make(chan string, 1) // надо использовать буферизированный канал
	go func(out chan<- string) {   // асинхронный запрос
		time.Sleep(1 * time.Second)
		fmt.Println("async operation ready, return user")
		out <- "user Eugene"
	}(result) // возвращает канал
	return result
}

func getPage() {
	// параллельный запуск:
	resultCh := getComments()
	userCh := getUser()

	time.Sleep(1 * time.Second)
	fmt.Println("get related articles")

	// забираем значение из канала:
	commentsData := <-resultCh // через 2 сек
	userData := <-userCh       // через 1 сек

	fmt.Println("main goroutine:", commentsData, ",", userData)
}

func main() {
	getPage()
}
