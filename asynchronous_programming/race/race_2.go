package main

import (
	"fmt"
	"sync" // содержит Mutex()
)

func main() {
	// все хорошо без GOMAXPROCS() при запуске
	var counters = map[int]int{}
	mu := &sync.Mutex{} // лучше создавать как ссылку
	// sm := sync.Map{} - безопасная мапа с которой можно работать из разных горутин
	for i := 0; i < 5; i++ {
		go func(counters map[int]int, th int, mu *sync.Mutex) {
			for j := 0; j < 5; j++ {
				mu.Lock() // другие горутины стоят в ожидании
				{
					counters[th*10+j]++
				}
				mu.Unlock()
			}
		}(counters, i, mu)
	}
	fmt.Scanln()
	mu.Lock()                                // для того, чтобы не ругался 'Race Detector'
	fmt.Println("counters result", counters) // можно не оборачивать в { }
	mu.Unlock()
} // все хорошо без GOMAXPROCS()
