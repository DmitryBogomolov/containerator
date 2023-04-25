/*
Program manage_container_server is an example of http server
that manages several container projects.

TODO:
  - add command status popups
*/
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/DmitryBogomolov/containerator/examples/manage_container_server/logger"
)

const defaultPort = 4001

func validateWorkspace(workspace string) error {
	cwd, _ := os.Getwd()
	stat, err := os.Stat(workspace)
	if err == nil && !stat.IsDir() {
		err = fmt.Errorf("'%s' is not a directory", workspace)
	}
	if err != nil {
		return err
	}
	abspath, _ := filepath.Abs(workspace)
	relpath, err := filepath.Rel(cwd, abspath)
	if err != nil {
		return err
	}
	if strings.HasPrefix(relpath, "..") {
		return fmt.Errorf("'%s' is outside working directory", workspace)
	}
	logger.Printf("workspace: %s\n", workspace)
	return nil
}

func runServer(port int, handler http.Handler) error {
	ch := make(chan error)
	go func() {
		ch <- http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
	}()
	logger.Printf("listening %d...\n", port)
	return <-ch
}

func run() error {
	var port int
	var workspace string
	flag.IntVar(&port, "port", defaultPort, "port")
	flag.StringVar(&workspace, "workspace", "", "path to workspace")
	flag.Parse()

	workspace, _ = filepath.Abs(workspace)
	err := validateWorkspace(workspace)
	if err != nil {
		return err
	}
	handler, err := setupServerHandler(workspace)
	if err != nil {
		return err
	}
	return runServer(port, handler)
}

func main() {
	if err := run(); err != nil {
		logger.Fatalf("%+v\n", err)
	}
}
