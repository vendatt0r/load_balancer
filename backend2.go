package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Response from backend 2")
	})

	fmt.Println("Backend 2 running on port 9002")
	http.ListenAndServe(":9002", nil)
}
