package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, h *http.Request) {
		fmt.Fprintln(w, "heelo")
	})

	http.ListenAndServe(":8080", nil)
}
