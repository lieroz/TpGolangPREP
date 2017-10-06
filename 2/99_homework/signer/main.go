package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"runtime"
	"sort"
	"strings"
)

const base = 10

func ExecutePipeline(jobs ...job) {

}

var m sync.Mutex // CRUTCH

func SingleHash(in, out chan interface{}) {
	data := <-out
	result := strconv.FormatInt(int64(data.(int)), base)
	ch := make(chan string, 1)
	go func(ch chan string) {
		ch <- DataSignerCrc32(<-ch)
	}(ch)
	ch <- result
	m.Lock()
	md5 := DataSignerMd5(result)
	m.Unlock()
	crc32 := <-ch
	in <- crc32 + "~" + DataSignerCrc32(md5)
}

func MultiHash(in, out chan interface{}) {
	singleHash := (<-in).(string)
	var chs [6]chan string
	var result string
	for i := 0; i < 6; i++ {
		chs[i] = make(chan string)
		go func(ch chan string, i int) {
			ch <- DataSignerCrc32(strconv.FormatInt(int64(i), base) + singleHash)
		}(chs[i], i)
	}
	for i := 0; i < 6; i++ {
		result += <-chs[i]
	}
	out <- result
}

func CombineResults(in, out chan interface{}) {
	var counter int
	result := make([]string, 0)
	for i := range out {
		counter++
		result = append(result, i.(string))
		if counter >= 7 {
			close(out)
		}
	}
	sort.Strings(result)
	in <- strings.Join(result[:], "_")
}

func main() {
	start := time.Now()
	in := make(chan interface{}, 1)
	out := make(chan interface{})
	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	for _, i := range inputData {
		go SingleHash(in, out)
		go MultiHash(in, out)
		out <- i
		runtime.Gosched() // CRUTCH
	}
	CombineResults(in, out)
	fmt.Println(<-in)
	fmt.Println(time.Since(start))
}
