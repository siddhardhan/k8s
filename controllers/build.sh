#!/bin/sh
mkdir -p src
rm -f bin/pod-annotator
go get 
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"'  -o bin/pod-annotator
echo "clean up cache..."
rm -rf src