build:
	CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -o bin/vpub main.go
