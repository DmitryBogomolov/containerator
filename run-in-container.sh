#!/bin/bash

src="$PWD"
dst="/go/${src//"$GOPATH/"}"
docker run -it  -v "$src":"$dst" -w "$dst" golang:1.11.0-stretch
