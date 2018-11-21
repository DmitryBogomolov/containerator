/*
Program manage_container_server is an example of http server
that manages several container projects.

TODO:
 * root page with basic description
 * page with list of all projects
 * button to deploy project
 * button to remove project
 * shows tags for image

*/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const defaultPort = 4001

func checkHTTPMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return false
	}
	return true
}

func setupServer() (http.Handler, error) {
	cache := newProjectsCache()
	rootHandler, err := newRootHandler(cache)
	if err != nil {
		return nil, err
	}
	commandHandler, err := newCommandHandler(cache)
	if err != nil {
		return nil, err
	}

	server := http.NewServeMux()
	server.Handle("/", rootHandler)
	server.Handle("/manage/", commandHandler)
	return server, nil
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

	handler, err := setupServer()
	if err != nil {
		log.Fatalf("%+v\n", err)
		os.Exit(1)
	}
	go runServer(port, handler, ch)
	log.Printf("Listening %d...", port)

	err = <-ch
	if err != nil {
		log.Fatalf("%+v\n", err)
		os.Exit(1)
	}
}
