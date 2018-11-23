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
	"encoding/json"
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
	server.Handle("/manage/", commandHandler)
	server.Handle("/static/", http.NotFoundHandler())
	server.HandleFunc("/static/index.js", func(w http.ResponseWriter, r *http.Request) {
		if !checkHTTPMethod(http.MethodGet, w, r) {
			return
		}
		http.ServeFile(w, r, "static/index.js")
	})
	server.HandleFunc("/api/projects", func(w http.ResponseWriter, r *http.Request) {
		if !checkHTTPMethod(http.MethodGet, w, r) {
			return
		}
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(cache.Projects); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !checkHTTPMethod(http.MethodGet, w, r) {
			return
		}
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		rootHandler.ServeHTTP(w, r)
	})
	return server, nil
}

type errorChan chan error

func runServer(port int, handler http.Handler, ch errorChan) {
	ch <- http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}

func run() error {
	var port int
	flag.IntVar(&port, "port", defaultPort, "port")
	flag.Parse()

	ch := make(errorChan)

	handler, err := setupServer()
	if err != nil {
		return err
	}

	go runServer(port, handler, ch)
	log.Printf("Listening %d...", port)

	return <-ch
}

func main() {
	err := run()
	if err != nil {
		log.Fatalf("%+v\n", err)
		os.Exit(1)
	}
}
