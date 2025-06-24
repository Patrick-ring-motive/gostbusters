go mod tidy

go fmt ./...

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o main