/*
Program manage_container_server is an example of http server
that manages several container projects.

TODO:
 * root page with basic description
 * page with list of all projects
 * button to deploy project
 * button to remove project
 * show tags for image
 * show responses in "tray"

*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/docker/docker/client"
)

const defaultPort = 4001

func checkMethod(method string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func checkPath(predicate func(p string) bool, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !predicate(r.URL.Path) {
			http.NotFound(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func indexScriptHandler() http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.js")
	})
	return checkMethod(http.MethodGet, handler)
}

func apiManageHandler(cache *projectsCache) http.Handler {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Panicln(err)
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.Replace(r.URL.Path, "/manage/", "", 1)
		item := cache.get(name)
		if item == nil {
			http.Error(w, fmt.Sprintf("'%s' is not found\n", name), http.StatusNotFound)
			return
		}
		body, err := invokeManage(cli, item.ConfigPath, r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error: %v\n", err), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, body)
	})
	return checkMethod(http.MethodPost, handler)
}

func apiTagsHandler(cache *projectsCache) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.Replace(r.URL.Path, "/api/tags/", "", 1)
		tags := []string{"3", "2", "1", name}
		data, err := json.Marshal(tags)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(data)
		w.Write([]byte("\n"))
	})
	return checkMethod(http.MethodGet, handler)
}

func createManageHandler(handler http.Handler) http.Handler {
	return checkMethod(http.MethodPost, handler)
}

func isRootPath(path string) bool {
	return path == "/"
}

func rootPageHandler(cache *projectsCache) http.Handler {
	tmpl := template.Must(template.ParseFiles("page.html"))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, cache)
		if err != nil {
			log.Printf("template error: %+v\n", err)
		}
	})
	return checkMethod(http.MethodGet, checkPath(isRootPath, handler))
}

func setupServer() (http.Handler, error) {
	cache := newProjectsCache()
	server := http.NewServeMux()
	server.Handle("/static/index.js", indexScriptHandler())
	server.Handle("/api/manage/", apiManageHandler(cache))
	server.Handle("/api/tags/", apiTagsHandler(cache))
	server.Handle("/", rootPageHandler(cache))
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
