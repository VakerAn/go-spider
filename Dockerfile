#FROM registry.srv.local/cn-dev-team/golang:1.20.4-alpine
#FROM scratch
#COPY --chown=33 . /app
#WORKDIR /app
FROM golang:latest

COPY --chown=33 . /app
WORKDIR /app

