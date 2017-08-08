deps:
	go get -u github.com/joho/godotenv
	go get -u github.com/aws/aws-sdk-go

build: test fmt
	go build

run:
	go run *.go

test:
	go test

fmt:
	go fmt

travis: deps fmt test build
