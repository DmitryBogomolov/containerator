/*
Program manage_container_server is an example of http server
that manages several container projects.

TODO:
 * refresh projects on 'not-found' error
 * add command status popups
*/
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const defaultPort = 4001

type errorChan chan error

func runServer(port int, handler http.Handler) error {
	ch := make(chan error)
	go func() {
		ch <- http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
	}()
	logger.Printf("Listening %d...\n", port)
	return <-ch
}

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
	p, err := filepath.Rel(cwd, abspath)
	if err != nil {
		return err
	}
	if strings.HasPrefix(p, "..") {
		return fmt.Errorf("'%s' is outside working directory", workspace)
	}
	return nil
}

func run() error {
	var port int
	var workspace string
	flag.IntVar(&port, "port", defaultPort, "port")
	flag.StringVar(&workspace, "workspace", ".sandbox", "path to workspace")
	flag.Parse()

	err := validateWorkspace(workspace)
	if err != nil {
		return err
	}
	logger.Printf("Workspace: %s\n", workspace)

	handler, err := setupServer(workspace)
	if err != nil {
		return err
	}

	return runServer(port, handler)
}

func main() {
	err := run()
	if err != nil {
		logger.Fatalf("%+v\n", err)
		os.Exit(1)
	}
}
