package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func uniq(input io.Reader, output io.Writer) error {
	in := bufio.NewScanner(input)
	var prev string
	for in.Scan() {
		txt := in.Text()
		if txt == prev {
			continue
		}
		if txt < prev {
			return fmt.Errorf("file not sorted")
		}
		prev = txt
		fmt.Fprintln(output, txt)
	}
	return nil
}

// запуск тестов: go test
// запуск тестов с описанием каждого: go test -v
// запуск всех тестов в текущем пакете и во всех подпакетов : go test -v ./...
func main() {
	err := uniq(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
