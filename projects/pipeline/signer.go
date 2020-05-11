package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type gorJob func(in, out chan interface{}, waiter *sync.WaitGroup)

const CountGorut = 100
const MultiHashCount = 6

var mux sync.Mutex

func Crc32(out chan interface{}, data string) {
	out <- DataSignerCrc32(data)
}

func GorSingleHash(in, out chan interface{}, waiter *sync.WaitGroup) {
	defer waiter.Done()
	for val := range in {
		data := fmt.Sprintf("%v", val)

		mux.Lock()
		md5 := DataSignerMd5(data)
		mux.Unlock()

		ch := make(chan interface{})
		chmd5 := make(chan interface{})

		go Crc32(ch, data)
		go Crc32(chmd5, md5)

		out <- fmt.Sprintf("%v", <-ch) + "~" + fmt.Sprintf("%v", <-chmd5)
	}
}

func GorMultHash(in, out chan interface{}, waiter *sync.WaitGroup) {
	defer waiter.Done()
	for val := range in {
		data := fmt.Sprintf("%v", val)

		wg := &sync.WaitGroup{}
		wg.Add(MultiHashCount)
		res := make([]string, MultiHashCount)
		for i := 0; i < MultiHashCount; i++ {
			go func(id int, str string, wg *sync.WaitGroup) {
				defer wg.Done()
				res[id] = DataSignerCrc32(strconv.Itoa(id) + str)
			}(i, data, wg)
		}
		wg.Wait()
		out <- strings.Join(res, "")
	}
}

func JobRunGur(in, out chan interface{}, f gorJob, cntGur int) {
	waiter := &sync.WaitGroup{}
	waiter.Add(cntGur)
	for i := 0; i < cntGur; i++ {
		go f(in, out, waiter)
	}
	waiter.Wait()
}

func MultiHash(in, out chan interface{}) {
	JobRunGur(in, out, GorMultHash, CountGorut)
}

func SingleHash(in, out chan interface{}) {
	JobRunGur(in, out, GorSingleHash, CountGorut)
}

func CombineResults(in, out chan interface{}) {
	hashes := make([]string, 0)

	for val := range in {
		data := fmt.Sprintf("%v", val)
		hashes = append(hashes, data)
	}
	sort.Strings(hashes)
	res := strings.Join(hashes, "_")
	out <- res
}

func newJob(j job, in, out chan interface{}, waiter *sync.WaitGroup) {
	defer waiter.Done()
	j(in, out)
	close(out)
}

func ExecutePipeline(jobs ...job) {
	waiter := &sync.WaitGroup{}
	pOut := make(chan interface{})
	for _, j := range jobs {
		cOut := make(chan interface{})
		waiter.Add(1)
		go newJob(j, pOut, cOut, waiter)
		pOut = cOut
	}
	waiter.Wait()
}

func main() {
}
