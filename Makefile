linux:
	env GOOS=linux GOARCH=amd64 go build -o build/monitgo_linux main.go