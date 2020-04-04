package main

import "fmt"

// обычная функция
func doNothing() {
	fmt.Println("i'm regular function")
}

func main() {
	// анонимная функция
	func(in string) {
		fmt.Println("anon func out:", in)
	}("nobody")

	// присванивание анонимной функции в переменную
	printer := func(in string) {
		fmt.Println("printer outs:", in)
	}
	printer("as variable") // теперь printer имеет тип func(string)
	printer = func(in string) {
		fmt.Println("other func.:", in)
	}

	// определяем тип функции
	type strFuncType func(string)

	// функция принимает коллбек
	worker := func(callback strFuncType) {
		callback("as callback")
	}
	worker(printer)

	someTest := "test"

	// функция возвращает замыкание
	prefixer := func(prefix string) strFuncType {
		return func(in string) {
			fmt.Printf("[%s] %s %s\n", someTest, prefix, in)
		}
	}
	successLogger := prefixer("SUCCESS")
	successLogger("expected behaviour")
	successLogger("bad behaviour")
}
