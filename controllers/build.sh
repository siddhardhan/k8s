#!/bin/sh
go get 
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"'  -o bin/pod-annotator