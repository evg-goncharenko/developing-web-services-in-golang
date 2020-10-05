package main

import "fmt"

func main() {
	ch1 := make(chan int) // канал
	//ch1 := make(chan int, 1) - буферизированный канал
	go func(in chan int) { // горутина
		val := <-in // чтение; состояние - режим ожидания
		fmt.Println("GO: get from chan", val)
		fmt.Println("GO: after read from chan")
	}(ch1)

	ch1 <- 42 // запись
	//ch1 <- 100500 - deadlock!

	fmt.Println("MAIN: after put to chan")
	fmt.Scanln()
}
