all: linux arm

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/monitgo_linux main.go

arm:
	GOOS=linux GOARCH=arm go build -o build/monitgo_arm main.go
