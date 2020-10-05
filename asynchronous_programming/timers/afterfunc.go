package main

import (
	"fmt"
	"time" // работа со временем
)

func sayHello() {
	fmt.Println("Hello World")
}

func main() {
	timer := time.AfterFunc(1*time.Second, sayHello) // выполнить через 1 секунду

	fmt.Scanln()
	timer.Stop() // остановка таймера, функция не будет выполнена

	timer = time.NewTimer(2 * time.Second) // через 2 секунды таймер сработает
	t := <-timer.C                         // timer имеет канал, в котором запишется момент времени, когда таймер сработает
	fmt.Println("Timer", t)

	t = <-time.After(1 * time.Second) // через 1 секунду таймер сработает
	fmt.Println("Time after", t)

	// таймер нельзя остановить!
}
