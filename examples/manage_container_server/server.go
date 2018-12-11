package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/docker/docker/client"
)

const (
	apiManageRoute = "/api/manage/"
	apiInfoRoute   = "/api/info/"
)

func restrictMethod(method string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
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

func apiManageHandler(cache *projectsCache, cli interface{}) http.Handler {
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

func apiInfoHandler(cache *projectsCache, cli interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		item, err := getProject(r, apiInfoRoute, cache)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		data, err := getImageInfo(cli, item.ConfigPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sendJSON(data, w)
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

func setupServer(pathToWorkspace string) (http.Handler, error) {
	cache := &projectsCache{
		Workspace: pathToWorkspace,
	}
	cache.refresh()

	cli, err := client.NewEnvClient()
	if err != nil {
		log.Panicln(err)
	}

	server := http.NewServeMux()
	server.Handle("/static/index.js", onlyGet(indexScriptHandler()))
	server.Handle(apiManageRoute, onlyPost(apiManageHandler(cache, cli)))
	server.Handle(apiInfoRoute, onlyGet(apiInfoHandler(cache, cli)))
	server.Handle("/", onlyRootPath(onlyGet(rootPageHandler(cache))))
	return server, nil
}
