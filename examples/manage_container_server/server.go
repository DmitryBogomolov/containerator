package main

import (
	_ "embed"
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	"github.com/DmitryBogomolov/containerator/examples/manage_container_server/logger"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
)

//go:embed static/index.js
var indexContent string

//go:embed static/page.html
var pageContent string
var pageTemplate = template.Must(template.New("/").Parse(pageContent))

func makeIndexScriptHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/javascript")
		w.Write([]byte(indexContent))
	})
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

func makeAPIManageHandler(cache *projectsCache, cli any) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		item, err := cache.get(vars["name"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		ret, err := invokeManage(cli, item.configPath, r)
		if err != nil {
			if os.IsNotExist(err) {
				cache.refresh()
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sendJSON(ret, w)
	})
}

func makeAPIInfoHandler(cache *projectsCache, cli any) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		item, err := cache.get(vars["name"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		data, err := getImageInfo(cli, item.configPath)
		if err != nil {
			if os.IsNotExist(err) {
				cache.refresh()
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sendJSON(data, w)
	})
}

func makeRootPageHandler(cache *projectsCache) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := pageTemplate.Execute(w, cache)
		if err != nil {
			logger.Printf("template error: %+v\n", err)
		}
	})
}

func attachLoggerToHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%-5s %s\n", r.Method, r.RequestURI)
		h.ServeHTTP(w, r)
	})
}

func setupServerHandler(pathToWorkspace string) (http.Handler, error) {
	cache := newProjectsCache(pathToWorkspace)

	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	server := mux.NewRouter()

	server.NewRoute().
		Path("/static/index.js").Methods(http.MethodGet).Handler(makeIndexScriptHandler())
	server.NewRoute().
		Path("/api/manage/{name}").Methods(http.MethodPost).Handler(makeAPIManageHandler(cache, cli))
	server.NewRoute().
		Path("/api/info/{name}").Methods(http.MethodGet).Handler(makeAPIInfoHandler(cache, cli))
	server.NewRoute().
		Path("/").Methods(http.MethodGet).Handler(makeRootPageHandler(cache))

	return attachLoggerToHandler(server), nil
}
