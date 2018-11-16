/*
Program manage_container_server is an example of http server
that manages several container projects.

TODO:
 * root page with basic description
 * page with list of all projects
 * button to deploy project
 * button to remove project
 * shows tags for image

*/
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/DmitryBogomolov/containerator"
	"github.com/docker/docker/client"

	"github.com/DmitryBogomolov/containerator/manage"
)

const defaultPort = 4001

var errBadURL = errors.New("bad url")
var errNoProject = errors.New("no project")

type commandHandler struct {
	cli     client.CommonAPIClient
	workDir string
}

var commandPattern = regexp.MustCompile("^/([\\w-]+)/([\\w-]+)")

func parseCommandURL(path string) (name string, cmd string, err error) {
	parts := commandPattern.FindStringSubmatch(path)
	if len(parts) == 0 {
		err = errBadURL
		return
	}
	name = parts[1]
	cmd = parts[2]
	return
}

func (h *commandHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name, command, err := parseCommandURL(r.URL.Path)
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
		return
	}
	configPath, err := findTarget(h.workDir, name)
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusNotFound)
		return
	}
	if command != "manage" {
		http.Error(w, "Error: bad command", http.StatusBadRequest)
		return
	}
	body, err := invokeManage(h.cli, configPath, r)
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, body)
}

func findTarget(workDir string, name string) (string, error) {
	items, err := filepath.Glob(filepath.Join(workDir, name, "*.yaml"))
	if err == nil && len(items) == 0 {
		err = errNoProject
	}
	if err != nil {
		return "", err
	}
	return items[0], nil
}

func parseBool(value string) bool {
	ret, _ := strconv.ParseBool(value)
	return ret
}

func invokeManage(cli client.CommonAPIClient, configPath string, r *http.Request) (string, error) {
	err := r.ParseForm()
	if err != nil {
		return "", err
	}
	config, err := manage.ReadConfig(configPath)
	if err != nil {
		return "", err
	}
	options := &manage.Options{
		Mode:   r.FormValue("mode"),
		Tag:    r.FormValue("tag"),
		Remove: parseBool(r.FormValue("remove")),
		Force:  parseBool(r.FormValue("force")),
	}
	cont, err := manage.Manage(cli, config, options)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Container: %s", containerator.GetContainerName(cont)), nil
}

func setupServer() (http.Handler, error) {
	workDir, _ := os.Getwd()
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	server := http.NewServeMux()
	server.Handle("/", &commandHandler{cli: cli, workDir: workDir})
	return server, nil
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

	handler, err := setupServer()
	if err != nil {
		log.Fatalf("%+v\n", err)
		os.Exit(1)
	}
	go runServer(port, handler, ch)
	log.Printf("Listening %d...", port)

	err = <-ch
	if err != nil {
		log.Fatalf("%+v\n", err)
		os.Exit(1)
	}
}
