package main

import (
	"fmt"
	"runtime"
)

func main() {
	// необходимо сделать какой-нибудь мониторинг
	runtime.GOMAXPROCS(1) // без этого ограничения будут ошибки: "concurrent map writes"
	var counters = map[int]int{}
	for i := 0; i < 5; i++ {
		go func(counters map[int]int, th int) { // запускаем 5 горутин
			for j := 0; j < 5; j++ {
				counters[th*10+j]++
			}
		}(counters, i)
	}
	fmt.Scanln()
	fmt.Println("counters result", counters)
}
