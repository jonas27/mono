FROM golang as builder
COPY go.mod go.sum cmd/main.go ./
RUN export GOPATH="" && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /metric-test .

RUN useradd scratchuser

FROM alpine
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs
COPY --from=builder /metric-test /
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/passwd /etc/passwd
USER scratchuser
CMD ["/metric-test"]