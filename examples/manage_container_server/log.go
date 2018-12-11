package main

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
