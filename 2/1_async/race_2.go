package main

import (
	"fmt"
	"sync"
)

var counters = map[int]int{}

func main() {
	mu := &sync.Mutex{}
	for i := 0; i < 5; i++ {
		go func(th int, protector *sync.Mutex) {
			for j := 0; j < 5; j++ {
				protector.Lock()
				counters[th*10+j]++
				protector.Unlock()
			}
		}(i, mu)
	}

	fmt.Scanln()
	fmt.Println("total result", counters)
}
