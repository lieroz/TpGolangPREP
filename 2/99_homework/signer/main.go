package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"sort"
	"strings"
)

const (
	th   = 6
	base = 10
)

func ExecutePipeline(jobs ...job) {

}

var m sync.Mutex

func SingleHash(in, out chan interface{}) {
	data := <-out
	result := strconv.FormatInt(int64(data.(int)), base)
	ch := make(chan string)
	go func(ch chan string) {
		ch <- DataSignerCrc32(<-ch)
	}(ch)
	ch <- result
	m.Lock()
	md5 := DataSignerMd5(result)
	m.Unlock()
	crc32md5 := DataSignerCrc32(md5)
	crc32 := <-ch
	in <- crc32 + "~" + crc32md5
}

func MultiHash(in, out chan interface{}) {
	singleHash := (<-in).(string)
	var chs [th]chan string
	var resultHash string
	for i := 0; i < th; i++ {
		chs[i] = make(chan string)
		go func(ch chan string, i int) {
			ch <- DataSignerCrc32(strconv.FormatInt(int64(i), base) + singleHash)
		}(chs[i], i)
	}
	for i := 0; i < th; i++ {
		resultHash += <-chs[i]
	}
	out <- resultHash
}

func CombineResults(in, out chan interface{}) {
	result := make([]string, 0)
	for hash := range out {
		result = append(result, hash.(string))
	}
	sort.Strings(result)
	in <- strings.Join(result[:], "_")
}

func main() {
	start := time.Now()
	in := make(chan interface{})
	out := make(chan interface{}, MaxInputDataLen)
	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	var wg sync.WaitGroup
	for _, i := range inputData {
		go SingleHash(in, out)
		wg.Add(1)
		go func() {
			defer wg.Done()
			MultiHash(in, out)
		}()
		out <- i
	}
	wg.Wait()
	close(out)
	go CombineResults(in, out)
	fmt.Println(<-in)
	fmt.Println(time.Since(start))
}
