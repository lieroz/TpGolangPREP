package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func uniq(input io.Reader) error {
	in := bufio.NewScanner(input)
	var prevText string
	for in.Scan() {
		text := in.Text()
		if prevText == text {
			continue
		}
		if text < prevText {
			return fmt.Errorf("file not sorted %s %s", text, prevText)
		}
		prevText = text
		fmt.Println(text)
	}
	return nil
}

func main() {
	err := uniq(os.Stdin)
	if err != nil {
		panic(err)
	}
}
