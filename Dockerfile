FROM ubuntu:22.04

COPY bin/rinha-de-backend /rinha-de-backend

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build api/main.go

COPY go.mod go.sum /rinha-de-backend


EXPOSE 9999
ENTRYPOINT  /rinha-de-backend
