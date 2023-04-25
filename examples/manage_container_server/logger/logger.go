package logger

import (
	"fmt"
	"log"
	"os"
)

type _LogWriter struct {
	file *os.File
}

func (w *_LogWriter) Write(b []byte) (int, error) {
	os.Stdout.Write(b)
	return w.file.Write(b)
}

const logFileName = "log.log"

var rootLogger = initRootLogger()

func initRootLogger() *log.Logger {
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	logWriter := _LogWriter{file}
	fmt.Fprintln(file, "----------------")
	return log.New(&logWriter, "", log.LstdFlags)
}

// Printf prints to the logger.
func Printf(format string, v ...any) {
	rootLogger.Printf(format, v...)
}

// Fatalf fatals to the logger.
func Fatalf(format string, v ...any) {
	rootLogger.Fatalf(format, v...)
}

// Panicf panics to the logger.
func Panicf(format string, v ...any) {
	rootLogger.Panicf(format, v...)
}
