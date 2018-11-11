package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const defaultPort = 4001

func handleCommand(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")[1:]
	if len(parts) != 2 {
		http.Error(w, "Bad url", http.StatusBadRequest)
		return
	}
	targetConfig := findTarget(parts[0])
	if targetConfig == "" {
		http.Error(w, "Bad name", http.StatusNotFound)
	}
	switch parts[1] {
	case "run":
		invokeRun(targetConfig, r)
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "OK\n")
	case "remove":
		invokeRemove(targetConfig)
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "OK\n")
	default:
		http.Error(w, "Bad command", http.StatusBadRequest)
	}
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

func invokeRun(config string, r *http.Request) {
}

func invokeRemove(config string) {
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
