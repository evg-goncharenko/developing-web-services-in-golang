package main

import (
	"fmt"
)

func main() {
	cancelCh := make(chan bool)
	dataCh := make(chan int)

	go func(cancelCh chan bool, dataCh chan int) {
		val := 0
		for {
			select { // проверка того, что пришла отмена или нет
			case <-cancelCh:
				fmt.Println("canceled")
				close(dataCh)
				return // асинхронная работа завершена
			case dataCh <- val:
				val++
			}
		}
	}(cancelCh, dataCh)

	for curVal := range dataCh {
		fmt.Println("read", curVal)
		if curVal > 3 {
			fmt.Println("send cancel")
			cancelCh <- true // будет ждать отработки в select
		}
	}

}
