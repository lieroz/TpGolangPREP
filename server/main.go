package main

import (
	"flag"
	"runtime"
)

func main() {
	port := flag.String("p", ":80", "a string")
	webRoot := flag.String("wr", "/var/www/html", "a string")
	numCpu := flag.Int("c", 0, "an int")
	workersCount := flag.Int("w", 4, "an int")
	flag.Parse()

	runtime.GOMAXPROCS(*numCpu)
	serv := NewServer(*port, *webRoot, *workersCount)
	serv.ListenAndServe()
}
