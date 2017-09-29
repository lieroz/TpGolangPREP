package main

import "fmt"

var counters = map[int]int{}

func main() {
	for i := 0; i < 5; i++ {
		go func(th int) {
			for j := 0; j < 5; j++ {
				counters[th*10+j]++
			}
		}(i)
	}

	fmt.Scanln()
	fmt.Println("money value", counters)
}
