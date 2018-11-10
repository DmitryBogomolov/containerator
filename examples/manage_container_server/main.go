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

func setupServer() *http.ServeMux {
	server := http.NewServeMux()
	server.HandleFunc("/run/", handleRun)
	server.HandleFunc("/remove/", handleRemove)
	return server
}

type errorChan chan error

func runServer(port int, handler http.Handler, ch errorChan) {
	ch <- http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}

func main() {
	var port int
	flag.IntVar(&port, "port", defaultPort, "port")
	flag.Parse()

	ch := make(errorChan)

	handler := setupServer()
	go runServer(port, handler, ch)
	log.Printf("Listening %d...", port)

	err := <-ch
	if err != nil {
		log.Fatalf("%+v\n", err)
		os.Exit(1)
	}
}
