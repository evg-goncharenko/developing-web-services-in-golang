package main

import (
	"fmt"
	"time"
)

var totalOperation int32 = 0

func increment() {
	totalOperation++
}

func main() {
	for i := 0; i < 1000; i++ { // запускаем 1000 горутин и делаем инкремент
		go increment()
	}
	time.Sleep(20 * time.Millisecond)
	// ожидается 1000, но по факту будет меньше
	fmt.Println("total operation = ", totalOperation)
}
