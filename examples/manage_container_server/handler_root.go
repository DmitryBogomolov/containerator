package main

import (
	"html/template"
	"net/http"
)

type rootHandler struct {
	cache    *projectsCache
	template *template.Template
}

func newRootHandler(cache *projectsCache) (*rootHandler, error) {
	tmpl := template.Must(template.ParseFiles("page.html"))
	return &rootHandler{
		cache:    cache,
		template: tmpl,
	}, nil
}

func (h *rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !checkHTTPMethod(http.MethodGet, w, r) {
		return
	}
	h.template.Execute(w, h.cache)
}
