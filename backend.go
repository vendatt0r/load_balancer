package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", "9001", "port to listen")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from backend at port %s", *port)
	})

	log.Printf("Backend listening on port %s", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
