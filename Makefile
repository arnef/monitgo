linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/monitgo_linux main.go