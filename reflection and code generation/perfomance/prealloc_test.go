// go test -bench . -benchmem prealloc_test.go
package main

import (
	"testing"
)

const iterNum = 1000

func BenchmarkEmptyAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data := make([]int, 0)
		for j := 0; j < iterNum; j++ {
			data = append(data, j)
		}
	}
}

func BenchmarkPreallocAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data := make([]int, 0, iterNum)
		for j := 0; j < iterNum; j++ {
			data = append(data, j)
		}
	}
}

/*
Запуск: (-x для генерации бинарного файла)
	go test -bench . -x -benchmem -cpuprofile=cpu.out -memprofile=mem.out -memprofilerate=1 prealloc_test.go

Запуск pprof: (после запуска консоли можно пользоваться командами (pprof))
	go tool pprof main.test cpu.out
	go tool pprof main.test mem.out

Запуск pprof, где больше тратилось места:
	go tool pprof -svg -inuse_space main.test.exe mem.out > mem_is.svg

Запуск pprof, где больше аллокаций:
	go tool pprof -svg -inuse_objects main.test.exe mem.out > mem_io.svg
	go tool pprof -svg main.test.exe cpu.out > cpu.svg

Запуск графической визуализации png:
	go tool pprof -png main.test.exe cpu.out > cpu.png

Запуск графической визуализации html:
	go tool pprof -web main.test cpu.out

Основные команды для pprof:
	(pprof) list .
	(pprof) list Benchmark*
	(pprof) inuse_space
	(pprof) list BenchmarkEmpty*
	(pprof) inuse_objects

	Дополнительные режимы:
	(pprof) alloc_objects
	(pprof) alloc_space
*/
