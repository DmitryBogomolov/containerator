package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const defaultPort = 4001

func handleRun(w http.ResponseWriter, r *http.Request) {
	log.Printf("run")
	fmt.Fprintf(w, "OK\n")
}

func handleRemove(w http.ResponseWriter, r *http.Request) {
	log.Printf("remove")
	fmt.Fprintf(w, "OK\n")
}

func setupServer() {
	http.HandleFunc("/run/", handleRun)
	http.HandleFunc("/remove/", handleRemove)
}

type errorChan chan error

func runServer(port int, ch errorChan) {
	ch <- http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func main() {
	var port int
	flag.IntVar(&port, "port", defaultPort, "port")
	flag.Parse()

	ch := make(errorChan)

	setupServer()
	go runServer(port, ch)
	log.Printf("Listening %d...", port)

	err := <-ch
	if err != nil {
		log.Fatalf("%+v\n", err)
		os.Exit(1)
	}
}
