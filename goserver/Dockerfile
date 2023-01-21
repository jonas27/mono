FROM golang:1.19 as builder
COPY go.mod go.sum /build/
WORKDIR /build/
RUN go mod vendor

# copy your data into container
COPY cmd/main.go /build/cmd/
