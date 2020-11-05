TLS 1.3 doesnt let you set the CipherSuite, use 1.2

Compile Benchmarks to binary
joma@INX-0102 MINGW64 ~/repos/mono/ecc-haproxy/go-client/cmd (master)
$ go test -c

joma@INX-0102 MINGW64 ~/repos/mono/ecc-haproxy/go-client/cmd (master)
$ ./cmd.test.exe -test.bench=.

GOOS=linux GOARCH=amd64 go test -c