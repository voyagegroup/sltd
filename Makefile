deps:
	go get -u github.com/joho/godotenv
	go get -u github.com/aws/aws-sdk-go
	go get -u github.com/fsnotify/fsnotify

build: test fmt
	go build

run:
	go run *.go

test:
	go test

fmt:
	go fmt

travis: fmt test build
