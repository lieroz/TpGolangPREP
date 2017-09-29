package main

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	iterationsNum = 3
	goroutinesNum = 12
	quotaLimit    = 3
)

func startWorker(in int, waiter *sync.WaitGroup, quota chan struct{}) {
	quota <- struct{}{} // ratelim.go, берём свободный слот на работу
	defer waiter.Done()
	for j := 0; j < iterationsNum; j++ {
		fmt.Printf(formatWork(in, j))

		// ratelim.go, раскомментируйте и посмотрите на результат
		<-quota             // ratelim.go, возвращаем слот
		quota <- struct{}{} // ratelim.go, берём свободный слот на работу

		runtime.Gosched() // даём поработать другим горутинам
	}
	<-quota // ratelim.go, возвращаем слот
}

func main() {
	runtime.GOMAXPROCS(1) // попробуйте с 0 (все доступные) и 1
	wg := &sync.WaitGroup{}
	quota := make(chan struct{}, quotaLimit) // ratelim.go
	for i := 0; i < goroutinesNum; i++ {
		wg.Add(1)
		go startWorker(i, wg, quota)
	}
	time.Sleep(time.Millisecond)
	wg.Wait() // wait_2.go ожидаем, пока waiter.Done() не приведёт счетчик к 0
}

func formatWork(in, j int) string {
	return fmt.Sprintln(strings.Repeat("  ", in), "█",
		strings.Repeat("  ", goroutinesNum-in),
		"th", in,
		"iter", j, strings.Repeat("■", j))
}
