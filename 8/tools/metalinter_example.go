package main

import (
	"fmt"
)

// MyStruct - my structure
type MyStruct struct {
	userID   int
	DataJSON []byte
}

// TestError func
func TestError(isOk bool) error {
	if !isOk {
		fmt.Errorf("failed")
	}
	return nil
}

// Test fff
func Test() {
	flag := true
	result := TestError(flag)
	fmt.Printf("result is\n", result)
	fmt.Printf("%v is %v", flag)

	s := &MyStruct{}
	fmt.Println(s)
}
