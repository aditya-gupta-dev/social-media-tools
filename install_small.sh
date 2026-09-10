CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o smt main.go
mv smt ~/go/bin/
