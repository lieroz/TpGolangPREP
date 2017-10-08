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
	Th   = 6
	Base = 10
)

//func ExecutePipeline(jobs ...job) {
//	var channels [len(jobs)]interface{}
//	for i := 0; i < len(jobs); i++ {
//		channels[i] = make(chan interface{}, MaxInputDataLen)
//	}
//}

var m sync.Mutex

func singleHash(out chan interface{}, data string) {
	ch := make(chan string)
	go func(ch chan string) {
		ch <- DataSignerCrc32(<-ch)
	}(ch)
	ch <- data
	m.Lock()
	md5 := DataSignerMd5(data)
	m.Unlock()
	crc32md5 := DataSignerCrc32(md5)
	crc32 := <-ch
	out <- crc32 + "~" + crc32md5
}

func SingleHash(in, out chan interface{}) {
	var wg sync.WaitGroup
	for val := range in {
		data := strconv.FormatInt(int64(val.(int)), Base)
		wg.Add(1)
		go func() {
			defer wg.Done()
			singleHash(out, data)
		}()
	}
	wg.Wait()
	close(out)
}

func multiHash(out chan interface{}, data string) {
	var chs [Th]chan string
	var resultHash string
	for i := 0; i < Th; i++ {
		chs[i] = make(chan string)
		go func(ch chan string, i int) {
			ch <- DataSignerCrc32(strconv.FormatInt(int64(i), Base) + data)
		}(chs[i], i)
	}
	for i := 0; i < Th; i++ {
		resultHash += <-chs[i]
	}
	out <- resultHash
}

func MultiHash(in, out chan interface{}) {
	var wg sync.WaitGroup
	for i := range in {
		data := i.(string)
		wg.Add(1)
		go func() {
			defer wg.Done()
			multiHash(out, data)
		}()
	}
	wg.Wait()
	close(out)
}

func CombineResults(in, out chan interface{}) {
	result := make([]string, 0)
	for val := range in {
		result = append(result, val.(string))
	}
	sort.Strings(result)
	out <- strings.Join(result[:], "_")
}

func main() {
	start := time.Now()
	ch1 := make(chan interface{}, MaxInputDataLen)
	ch2 := make(chan interface{}, MaxInputDataLen)
	ch3 := make(chan interface{}, MaxInputDataLen)
	ch4 := make(chan interface{}, MaxInputDataLen)
	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	go func() {
		for _, i := range inputData {
			ch1 <- i
		}
		close(ch1)
	}()
	go SingleHash(ch1, ch2)
	go MultiHash(ch2, ch3)
	go CombineResults(ch3, ch4)
	fmt.Println(<-ch4)
	fmt.Println(time.Since(start))
}
