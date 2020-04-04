package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	iterationsNum = 7
	goroutinesNum = 5
)

func startWorker(in int, waiter *sync.WaitGroup) {
	defer waiter.Done() // декрементим wait группу (уменьшаем счетчик на 1)
	for j := 0; j < iterationsNum; j++ {
		fmt.Printf(formatWork(in, j))
		time.Sleep(time.Millisecond)
	}
}

func main() {
	wg := &sync.WaitGroup{} // инициализируем wait группу
	for i := 0; i < goroutinesNum; i++ {
		// wg.Add надо вызывать в той горутине, которая порождает воркеров
		// иначе другая горутина может не успеть запуститься и выполнится Wait
		wg.Add(1) // инкрементим wait группу (увеличиваем счетчик на 1)
		go startWorker(i, wg)
	}
	time.Sleep(time.Millisecond)
	wg.Wait() // ожидаем, пока waiter.Done() не приведёт счетчик к 0

	fmt.Println("end")

}

func formatWork(in, j int) string {
	return fmt.Sprintln(strings.Repeat("  ", in), "█",
		strings.Repeat("  ", goroutinesNum-in),
		"th", in,
		"iter", j, strings.Repeat("■", j))
}
