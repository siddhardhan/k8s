#!/bin/sh
echo "Building the code..."
GOOS=linux go build -o ./app .
eval $(minikube docker-env)
echo "Building image..."
docker build -t in-cluster .