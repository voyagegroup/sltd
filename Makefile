deps:
	go get -u github.com/joho/godotenv
	go get -u github.com/aws/aws-sdk-go
	go get -u github.com/fsnotify/fsnotify
	go get -u github.com/mitchellh/gox
	go get -u github.com/tcnksm/ghr

build: test fmt
	go build

run:
	go run *.go

test:
	go test

fmt:
	go fmt

travis: fmt test build

dist: dist/clean dist/build dist/pack dist/upload

dist/clean:
	mkdir -p pkg/ dist/
	rm -rf pkg/*
	rm -rf dist/*

dist/build: XC_ARCH=386 amd64
dist/build: XC_OS=linux darwin windows
dist/build: test fmt
	gox \
	    -os="$(XC_OS)" \
	    -arch="$(XC_ARCH)" \
	    -output "pkg/{{.Dir}}_{{.OS}}_{{.Arch}}/{{.Dir}}_{{.OS}}_{{.Arch}}"

dist/pack:
	@for DIR in $$(ls pkg/ | grep -v dist/); do \
		zip -j dist/$${DIR}.zip pkg/$${DIR}/*; \
	done

dist/upload: RELEASE_VERSION=latest
dist/upload:
	ghr -u voyagegroup -r sltd -recreate "$(RELEASE_VERSION)" dist/
