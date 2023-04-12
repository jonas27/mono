From golang:1.20
WORKDIR /build
COPY go.sum go.mod ./
RUN go mod download
COPY main* ./
RUN go test -run TestInxAuth