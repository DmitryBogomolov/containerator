package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const defaultPort = 4001

var errBadCommand = errors.New("bad command")

func handleCommand(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")[1:]
	if len(parts) != 2 {
		http.Error(w, "Bad url", http.StatusBadRequest)
		return
	}
	targetConfig := findTarget(parts[0])
	if targetConfig == "" {
		http.Error(w, "Bad name", http.StatusNotFound)
		return
	}
	body := ""
	err := errBadCommand
	if parts[1] == "manage" {
		body, err = invokeManage(targetConfig, r)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, body)
}

func findTarget(name string) string {
	dir, _ := os.Getwd()
	items, _ := ioutil.ReadDir(dir)
	for _, item := range items {
		if item.IsDir() && item.Name() == name {
			return filepath.Join(dir, name, "config.yaml")
		}
	}
	return ""
}

func invokeManage(config string, r *http.Request) (string, error) {
	err := r.ParseForm()
	if err != nil {
		return "", err
	}
	mode := r.FormValue("mode")
	tag := r.FormValue("tag")
	force := r.FormValue("force")
	remove := r.FormValue("remove")
	return mode + tag + force + remove, nil
}

func setupServer() *http.ServeMux {
	server := http.NewServeMux()
	server.HandleFunc("/", handleCommand)
	return server
}

type errorChan chan error

func runServer(port int, handler http.Handler, ch errorChan) {
	ch <- http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}

func main() {
	var port int
	flag.IntVar(&port, "port", defaultPort, "port")
	flag.Parse()

	ch := make(errorChan)

	handler := setupServer()
	go runServer(port, handler, ch)
	log.Printf("Listening %d...", port)

	err := <-ch
	if err != nil {
		log.Fatalf("%+v\n", err)
		os.Exit(1)
	}
}
