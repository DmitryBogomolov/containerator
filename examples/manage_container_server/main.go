/*
Program manage_container_server is an example of http server
that manages several container projects.

TODO:
 * refresh projects on time interval
 * refresh projects on 'not-found' error
 * add command status popups
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
const apiManageRoute = "/api/manage/"
const apiTagsRoute = "/api/tags/"

func restrictMethod(method string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func onlyGet(handler http.Handler) http.Handler {
	return restrictMethod(http.MethodGet, handler)
}

func onlyPost(handler http.Handler) http.Handler {
	return restrictMethod(http.MethodPost, handler)
}

func restrictPath(predicate func(p string) bool, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !predicate(r.URL.Path) {
			http.NotFound(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func isRootPath(path string) bool {
	return path == "/"
}

func onlyRootPath(handler http.Handler) http.Handler {
	return restrictPath(isRootPath, handler)
}

func indexScriptHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.js")
	})
}

func getProject(r *http.Request, prefix string, cache *projectsCache) (*projectItem, error) {
	name := strings.TrimPrefix(r.URL.Path, prefix)
	item := cache.get(name)
	if item == nil {
		return nil, fmt.Errorf("project '%s' is not found", name)
	}
	return item, nil
}

func sendJSON(value interface{}, w http.ResponseWriter) {
	data, err := json.Marshal(value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
	w.Write([]byte("\n"))
}

func apiManageHandler(cache *projectsCache) http.Handler {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Panicln(err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		item, err := getProject(r, apiManageRoute, cache)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		ret, err := invokeManage(cli, item.ConfigPath, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sendJSON(ret, w)
	})
}

func apiTagsHandler(cache *projectsCache) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		item, err := getProject(r, apiTagsRoute, cache)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		tags := []string{"3", "2", "1", item.Name}
		sendJSON(tags, w)
	})
}

func rootPageHandler(cache *projectsCache) http.Handler {
	tmpl := template.Must(template.ParseFiles("page.html"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, cache)
		if err != nil {
			log.Printf("template error: %+v\n", err)
		}
	})
}

func setupServer() (http.Handler, error) {
	cache := newProjectsCache()
	server := http.NewServeMux()
	server.Handle("/static/index.js", onlyGet(indexScriptHandler()))
	server.Handle(apiManageRoute, onlyPost(apiManageHandler(cache)))
	server.Handle(apiTagsRoute, onlyGet(apiTagsHandler(cache)))
	server.Handle("/", onlyRootPath(onlyGet(rootPageHandler(cache))))
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
