# syntax=docker/dockerfile:1

# FROM golang:1.16.13-alpine3.14
FROM golang:latest

WORKDIR /app

COPY go.* ./

RUN go mod download

# Copy local code to the container image.
COPY ./src/* ./

RUN go build -o /main

CMD [ "/main" ]
