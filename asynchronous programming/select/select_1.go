package main

import (
	"fmt"
)

func main() {
	ch1 := make(chan int, 1) // буферизированный канал 1
	ch2 := make(chan int)    // канал 2

	select { // спец. оператор, который позволяет выбрать один из доступных ассинхронных операций
	case val := <-ch1: // попытка получить данные из канала 1
		fmt.Println("ch1 val", val)
	case ch2 <- 1: // попытка записать данные в канал 2
		fmt.Println("put val to ch2")
	default:
		fmt.Println("default case")
	}
}
