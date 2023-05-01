From golang:1.20 as builder
WORKDIR /build
COPY go.sum go.mod ./
RUN go mod download
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /kustomize-check ./

FROM scratch 
COPY --from=builder /kustomize-check /kustomize-check
ENTRYPOINT [ "/kustomize-check" ]
