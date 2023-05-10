FROM golang:1.20-bullseye

ENV GO111MODULE=on

WORKDIR /app/src/cafe

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

EXPOSE 8080
