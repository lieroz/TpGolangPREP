package main

import (
	"fmt"
)

func main() {

	cancelCh := make(chan struct{})
	dataCh := make(chan int)

	go func(cancelCh chan struct{}, dataCh chan int) {
		val := 0
		for {
			select {
			case <-cancelCh:
				return
			case dataCh <- val:
				fmt.Println("put", val)
				val++
				// ничего не делаем, сама операция - отправка в канал
				// default:
				// 	fmt.Println("default", val)
				// 	// переход в след итерацию цикла
			}
		}
	}(cancelCh, dataCh)

	for curVal := range dataCh {
		if curVal > 10 {
			cancelCh <- struct{}{}
			break
		}
		fmt.Println("read", curVal)
	}

}
