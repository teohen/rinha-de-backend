#!/bin/bash
export LOCAL_ENV=true
export HTTP_PORT=9999
export DB_URL="localhost://postgres@postgres:5432/pessoas?=sslmode=disable"
export TRACE_ENABLED=true

CGO_ENABLED=0 go build  -o ./bin/rinha-de-backend .

./bin/rinha-de-backend

