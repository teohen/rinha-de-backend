# syntax=docker/dockerfile:1
FROM golang:1.21 as base

WORKDIR /app
COPY . /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build api/main.go
FROM ubuntu:22.04

COPY go.mod go.sum ./
WORKDIR /app

COPY --from=base /app/main ./

EXPOSE 80
ENV GIN_MODE release
CMD [ "./main" ]
