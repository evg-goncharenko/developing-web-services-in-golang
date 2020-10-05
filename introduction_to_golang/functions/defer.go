package main

import "fmt"

func getSomeVars() string {
	fmt.Println("getSomeVars execution")
	return "getSomeVars result"
}

func main() {
	// выполняется по завершению функции
	defer fmt.Println("After work")
	defer func() {
		fmt.Println(getSomeVars())
	}()
	fmt.Println("Some userful work")

	/*
		Some userful work
		getSomeVars execution
		getSomeVars result
		After work
	*/
}
