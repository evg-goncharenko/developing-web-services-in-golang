package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second)
	i := 0
	for tickTime := range ticker.C { // проходимся по каналу
		i++
		fmt.Println("step", i, "time", tickTime)
		if i >= 5 {
			// надо останавливать, иначе утечка памяти
			ticker.Stop()
			break
		}
	}
	fmt.Println("total", i)

	// не может быть остановлен и собран сборщиком мусора
	// используем если должен работать вечно, тк нет Stop()
	c := time.Tick(time.Second) // time.Tick возвращает канал
	i = 0
	for tickTime := range c {
		i++
		fmt.Println("step", i, "time", tickTime)
		if i >= 5 {
			break
		}
	}

}
