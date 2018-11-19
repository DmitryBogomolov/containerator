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
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/DmitryBogomolov/containerator"
	"github.com/docker/docker/client"

	"github.com/DmitryBogomolov/containerator/manage"
)

const defaultPort = 4001

var errBadURL = errors.New("bad url")
var errNoProject = errors.New("no project")

type rootHandler struct {
	template *template.Template
	workDir  string
}

func newRootHandler(workDir string) (*rootHandler, error) {
	template, err := template.ParseFiles("page.html")
	if err != nil {
		return nil, err
	}
	return &rootHandler{
		template: template,
		workDir:  workDir,
	}, nil
}

func (h *rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	data := struct {
		Message string
	}{
		Message: "Hello World",
	}
	h.template.Execute(w, &data)
}

type commandHandler struct {
	cli     client.CommonAPIClient
	workDir string
}

func newCommandHandler(workDir string) (*commandHandler, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return &commandHandler{
		cli:     cli,
		workDir: workDir,
	}, nil
}

func (h *commandHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	name := strings.Replace(r.URL.Path, "/manage/", "", 1)
	configPath, err := findTarget(h.workDir, name)
	if err != nil {
		http.Error(w, fmt.Sprintf("'%s' is not found\n", name), http.StatusNotFound)
		return
	}
	body, err := invokeManage(h.cli, configPath, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v\n", err), http.StatusInternalServerError)
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
	options := &manage.Options{
		Mode:   r.PostFormValue("mode"),
		Tag:    r.PostFormValue("tag"),
		Remove: parseBool(r.PostFormValue("remove")),
		Force:  parseBool(r.PostFormValue("force")),
	}
	config, err := manage.ReadConfig(configPath)
	if err != nil {
		return "", err
	}
	cont, err := manage.Manage(cli, config, options)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Container: %s", containerator.GetContainerName(cont)), nil
}

func setupServer() (http.Handler, error) {
	workDir, _ := os.Getwd()
	rootHandler, err := newRootHandler(workDir)
	if err != nil {
		return nil, err
	}
	commandHandler, err := newCommandHandler(workDir)
	if err != nil {
		return nil, err
	}

	server := http.NewServeMux()
	server.Handle("/", rootHandler)
	server.Handle("/manage/", commandHandler)
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
