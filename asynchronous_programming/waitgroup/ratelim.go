package main

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	iterationsNumbers = 6
	goroutinesNumbers = 5 // всего воркеров - 5
	quotaLimit        = 2 // параллельно работающих воркеров - 2
)

func startWorkers(in int, wg *sync.WaitGroup, quotaCh chan struct{}) {
	quotaCh <- struct{}{} // канал с квотой, берём свободный слот, аналог семафора
	defer wg.Done()
	for j := 0; j < iterationsNumbers; j++ {
		fmt.Printf(formatWorks(in, j))

		if j%2 == 0 {
			<-quotaCh             // возвращаем слот
			quotaCh <- struct{}{} // берём слот
		}

		runtime.Gosched() // даём поработать другим горутинам
	}
	<-quotaCh // возвращаем слот
}

func main() {
	wg := &sync.WaitGroup{}
	quotaCh := make(chan struct{}, quotaLimit)
	for i := 0; i < goroutinesNumbers; i++ {
		wg.Add(1)
		go startWorkers(i, wg, quotaCh)
	}
	time.Sleep(time.Millisecond)
	wg.Wait()
}

func formatWorks(in, j int) string {
	return fmt.Sprintln(strings.Repeat("  ", in), "█",
		strings.Repeat("  ", goroutinesNumbers-in),
		"th", in,
		"iter", j, strings.Repeat("■", j))
}
