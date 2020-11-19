package main

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

/*
	this is a test to check that we really have a pipeline
	wrong behavior: accumulate the results of one function, and then send them to the next
	this doesn't allow you to run infinite tasks on the pipeline
	correct behavior: provide smooth flow
*/

func TestPipeline(t *testing.T) {
	var ok = true
	var recieved uint32
	freeFlowJobs := []job{
		job(func(in, out chan interface{}) {
			out <- 1
			time.Sleep(10 * time.Millisecond)
			currRecieved := atomic.LoadUint32(&recieved)
			/*
				if the values accumulate, then until the entire function works , they will not go further
				here we check that the counter has increased in the following function
				it means that the value got there before the current function worked
			*/
			if currRecieved == 0 {
				ok = false
			}
		}),
		job(func(in, out chan interface{}) {
			for _ = range in {
				atomic.AddUint32(&recieved, 1)
			}
		}),
	}
	ExecutePipeline(freeFlowJobs...)
	if !ok || recieved == 0 {
		t.Errorf("no value free flow - dont collect them")
	}
}

func TestSigner(t *testing.T) {
	testExpected := "1173136728138862632818075107442090076184424490584241521304_1696913515191343735512658979631549563179965036907783101867_27225454331033649287118297354036464389062965355426795162684_29568666068035183841425683795340791879727309630931025356555_3994492081516972096677631278379039212655368881548151736_4958044192186797981418233587017209679042592862002427381542_4958044192186797981418233587017209679042592862002427381542"
	testResult := "NOT_SET"
	/*
		this is a small protection against attempts not to call the calculation functions
		redefinition of functions on their own, which increment the local counter
		redefinition is possible because the function is declared as a variable in which the function lies
	*/
	var (
		DataSignerSalt         string = "" // the server will have a different value
		OverheatLockCounter    uint32
		OverheatUnlockCounter  uint32
		DataSignerMd5Counter   uint32
		DataSignerCrc32Counter uint32
	)
	OverheatLock = func() {
		atomic.AddUint32(&OverheatLockCounter, 1)
		for {
			if swapped := atomic.CompareAndSwapUint32(&dataSignerOverheat, 0, 1); !swapped {
				fmt.Println("OverheatLock happend")
				time.Sleep(time.Second)
			} else {
				break
			}
		}
	}
	OverheatUnlock = func() {
		atomic.AddUint32(&OverheatUnlockCounter, 1)
		for {
			if swapped := atomic.CompareAndSwapUint32(&dataSignerOverheat, 1, 0); !swapped {
				fmt.Println("OverheatUnlock happend")
				time.Sleep(time.Second)
			} else {
				break
			}
		}
	}
	DataSignerMd5 = func(data string) string {
		atomic.AddUint32(&DataSignerMd5Counter, 1)
		OverheatLock()
		defer OverheatUnlock()
		data += DataSignerSalt
		dataHash := fmt.Sprintf("%x", md5.Sum([]byte(data)))
		time.Sleep(10 * time.Millisecond)
		return dataHash
	}
	DataSignerCrc32 = func(data string) string {
		atomic.AddUint32(&DataSignerCrc32Counter, 1)
		data += DataSignerSalt
		crcH := crc32.ChecksumIEEE([]byte(data))
		dataHash := strconv.FormatUint(uint64(crcH), 10)
		time.Sleep(time.Second)
		return dataHash
	}

	inputData := []int{0, 1, 1, 2, 3, 5, 8}

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			dataRaw := <-in
			data, ok := dataRaw.(string)
			if !ok {
				t.Error("can't convert result data to string")
			}
			testResult = data
		}),
	}

	start := time.Now()

	ExecutePipeline(hashSignJobs...)

	end := time.Since(start)

	expectedTime := 3 * time.Second

	if testExpected != testResult {
		t.Errorf("results not match\nGot: %v\nExpected: %v", testResult, testExpected)
	}

	if end > expectedTime {
		t.Errorf("execition too long\nGot: %s\nExpected: <%s", end, time.Second*3)
	}

	if int(OverheatLockCounter) != len(inputData) ||
		int(OverheatUnlockCounter) != len(inputData) ||
		int(DataSignerMd5Counter) != len(inputData) ||
		int(DataSignerCrc32Counter) != len(inputData)*8 { // 8 because 2 in SingleHash and 6 in MultiHash
		t.Errorf("not enough hash-func calls")
	}
}
