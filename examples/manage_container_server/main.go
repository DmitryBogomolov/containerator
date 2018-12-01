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

func checkMethod(method string, handler http.Handler) http.Handler {
	check := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		handler.ServeHTTP(w, r)
	}
	return http.HandlerFunc(check)
}

func checkPath(predicate func(p string) bool, handler http.Handler) http.Handler {
	check := func(w http.ResponseWriter, r *http.Request) {
		if !predicate(r.URL.Path) {
			http.NotFound(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	}
	return http.HandlerFunc(check)
}

func createScriptHandler() http.Handler {
	handle := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.js")
	}
	return checkMethod(http.MethodGet, http.HandlerFunc(handle))
}

func createAPIHandler(cache *projectsCache) http.Handler {
	handle := func(w http.ResponseWriter, r *http.Request) {
		names := make([]string, len(cache.Projects))
		for i, p := range cache.Projects {
			names[i] = p.Name
		}
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(names); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	return checkMethod(http.MethodGet, http.HandlerFunc(handle))
}

func createManageHandler(handler http.Handler) http.Handler {
	return checkMethod(http.MethodPost, handler)
}

func createPageHandler(handler http.Handler) http.Handler {
	h := checkPath(func(path string) bool { return path == "/" }, handler)
	return checkMethod(http.MethodGet, h)
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
	server.Handle("/manage/", createManageHandler(commandHandler))
	server.Handle("/static/index.js", createScriptHandler())
	server.Handle("/api/projects", createAPIHandler(cache))
	server.Handle("/", createPageHandler(rootHandler))
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
