package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	"github.com/DmitryBogomolov/containerator/examples/manage_container_server/logger"
	"github.com/gorilla/mux"

	"github.com/docker/docker/client"
)

func indexScriptHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.js")
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

func apiManageHandler(cache *projectsCache, cli interface{}) http.Handler {
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

func apiInfoHandler(cache *projectsCache, cli interface{}) http.Handler {
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

func rootPageHandler(cache *projectsCache) http.Handler {
	tmpl := template.Must(template.ParseFiles("page.html"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, cache)
		if err != nil {
			logger.Printf("template error: %+v\n", err)
		}
	})
}

func wrapLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%-5s %s\n", r.Method, r.RequestURI)
		h.ServeHTTP(w, r)
	})
}

func setupServer(pathToWorkspace string) (http.Handler, error) {
	cache := newProjectsCache(pathToWorkspace)

	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	server := mux.NewRouter()

	server.NewRoute().
		Path("/static/index.js").Methods(http.MethodGet).Handler(indexScriptHandler())
	server.NewRoute().
		Path("/api/manage/{name}").Methods(http.MethodPost).Handler(apiManageHandler(cache, cli))
	server.NewRoute().
		Path("/api/info/{name}").Methods(http.MethodGet).Handler(apiInfoHandler(cache, cli))
	server.NewRoute().
		Path("/").Methods(http.MethodGet).Handler(rootPageHandler(cache))

	return wrapLogger(server), nil
}
