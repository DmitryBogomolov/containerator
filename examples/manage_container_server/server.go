package main

import (
	_ "embed"
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/DmitryBogomolov/containerator/examples/manage_container_server/logger"
	"github.com/DmitryBogomolov/containerator/examples/manage_container_server/registry"
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

func sendJSON(value any, w http.ResponseWriter) {
	data, err := json.Marshal(value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
	w.Write([]byte("\n"))
}

func makeAPIManageContainerHandler(registry *registry.Registry, cli any) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		targetName := mux.Vars(r)["name"]
		registry.Refresh()
		item, err := registry.GetItem(targetName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		data, err := invokeManage(cli, item.ConfigPath, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sendJSON(data, w)
	})
}

func makeAPIImageInfoHandler(registry *registry.Registry, cli any) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		targetName := mux.Vars(r)["name"]
		registry.Refresh()
		item, err := registry.GetItem(targetName)
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

func makeRootPageHandler(registry *registry.Registry) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := pageTemplate.Execute(w, registry.Items())
		if err != nil {
			logger.Printf("template error: %+v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
	registry := registry.New(pathToWorkspace)

	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	server := mux.NewRouter()

	server.NewRoute().
		Path("/static/index.js").
		Methods(http.MethodGet).
		Handler(makeIndexScriptHandler())
	server.NewRoute().
		Path("/api/manage-container/{name}").
		Methods(http.MethodPost).
		Handler(makeAPIManageContainerHandler(registry, cli))
	server.NewRoute().
		Path("/api/image-info/{name}").
		Methods(http.MethodGet).
		Handler(makeAPIImageInfoHandler(registry, cli))
	server.NewRoute().
		Path("/").
		Methods(http.MethodGet).
		Handler(makeRootPageHandler(registry))

	return attachLoggerToHandler(server), nil
}
