package main

import (
	"fmt"
)

func deferTest() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic happend FIRST:", err)
			err = fmt.Errorf("panic happend, %v", err)
		}
	}()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic happend SECOND:", err)
			panic("second panic")
		}
	}()
	fmt.Println("Some userful work")
	panic("something bad happend")
	return
}

// в своем коде panic не должно быть никогда
func main() {
	deferTest()
	return
}
