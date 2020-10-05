package main

import "fmt"

const pi = 3.141
const (
	hello = "Привет"
	e     = 2.718
)
const (
	zero  = iota // iota - автоинкремент для констант, = 0
	_            // пустая переменная, пропуск iota, = 1
	two          // = 2
	three        // = 3
)
const (
	_ = iota // пропускаем первое значение, = 0
	// KB 1 << (10 * 1) = 1024
	KB uint64 = 1 << (10 * iota) // iota = 1
	// MB 1 << (10 * 2) = 1048576
	MB // iota = 2
)
const (
	// нетипизированная константа
	year = 2017
	// типизированная константа
	yearTyped int = 2017
)

func main() {
	var month int32 = 13
	fmt.Println(month + year)

	// month + yearTyped - нельзя (mismatched types int32 and int)
}
