package main

import (
	"html/template"
	"log"
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
	err := h.template.Execute(w, h.cache)
	if err != nil {
		log.Printf("template error: %+v\n", err)
	}
}
