package main

import (
	"bufio"
	"fmt"
	"os"
)

// запуск: cat data.txt | go run main.go
func main() {
	// NewScanner - для построчного чтения
	in := bufio.NewScanner(os.Stdin)
	var prev string
	for in.Scan() {
		txt := in.Text()
		if txt == prev {
			continue
		}
		if txt < prev {
			panic("file not sorted")
		}
		prev = txt
		fmt.Println(txt)
	}
}
