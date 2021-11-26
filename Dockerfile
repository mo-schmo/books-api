# syntax=docker/dockerfile:1

FROM golang:latest

ENV PORT 8083

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod tidy

COPY . ./

EXPOSE $PORT

CMD [ "go", "run", "main.go" ]