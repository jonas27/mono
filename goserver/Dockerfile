FROM golang:1.19-alpine as builder
WORKDIR /go/src/goserver
COPY go.mod go.sum ./
RUN go mod download 
COPY cmd/main.go cmd/main.go
WORKDIR /go/src/goserver/cmd
