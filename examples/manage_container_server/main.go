/*
Program manage_container_server is an example of http server
that manages several container projects.

TODO:
 * refresh projects on time interval
 * refresh projects on 'not-found' error
 * add command status popups
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const defaultPort = 4001

type errorChan chan error

func runServer(port int, handler http.Handler, ch errorChan) {
	ch <- http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}

func run() error {
	var port int
	flag.IntVar(&port, "port", defaultPort, "port")
	flag.Parse()

	ch := make(errorChan)

	handler, err := setupServer()
	if err != nil {
		return err
	}

	go runServer(port, handler, ch)
	log.Printf("Listening %d...", port)

	return <-ch
}

func main() {
	err := run()
	if err != nil {
		log.Fatalf("%+v\n", err)
		os.Exit(1)
	}
}
