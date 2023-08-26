#!/bin/bash

GOOS=linux GOARCH=amd64 CGO-ENABLED=0 go build -o ./bin/rinha-de-backend .

image="teohen/rinha-de-backend"
tag="local"

docker build --no-cache=true -t "${image}:${tag}" .
