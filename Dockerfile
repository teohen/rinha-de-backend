FROM ubuntu:22.04

COPY bin/rinha-de-backend /rinha-de-backend

EXPOSE 9999
ENTRYPOINT  /rinha-de-backend
