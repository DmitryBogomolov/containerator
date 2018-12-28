package main

import (
	"fmt"
	"log"
	"os"
)

type logWriter struct {
	file *os.File
}

func (w *logWriter) Write(b []byte) (int, error) {
	os.Stdout.Write(b)
	return w.file.Write(b)
}

var logger = createLogger()

func createLogger() *log.Logger {
	file, err := os.OpenFile("log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	w := logWriter{file}
	fmt.Fprintln(file, "----------------")
	return log.New(&w, "", log.LstdFlags)
}
