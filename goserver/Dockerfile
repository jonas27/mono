FROM golang:1.19 as builder
COPY go.mod go.sum /build/
COPY cmd/main.go /build/cmd/
WORKDIR /build/
RUN go mod vendor
