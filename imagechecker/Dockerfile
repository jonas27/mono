FROM golang:1.19.4-alpine as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY pkg pkg

# Build the application
RUN go build -o /imagechecker ./pkg/

FROM scratch
ARG VERSION
ENV VERSION=$VERSION
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /imagechecker /imagechecker
ENTRYPOINT [ "/imagechecker" ]