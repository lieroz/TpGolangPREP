package main

import "fmt"

func main() {

	ch1 := make(chan int)

	go func(in <-chan int) {
		val := <-in
		fmt.Println(val)
		fmt.Println("after read from chan")
	}(ch1)

	ch1 <- 1

	// вызывает дедлок на небуферизированном канале
	fmt.Println(<-ch1)

	fmt.Println("after put to chan")

	fmt.Scanln()
}
