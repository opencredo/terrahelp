BUILDARGS ?= -mod=vendor

check:
	go vet $(BUILDARGS) ./...

test:
	go test $(BUILDARGS) -v ./...

build:
	go build $(BUILDARGS) -o ${GOPATH}/bin/terrahelp

dist:
	- GOOS=darwin GOARCH=amd64 go build $(BUILDARGS) -o=terrahelp-darwin-amd64
	- GOOS=linux GOARCH=amd64 go build $(BUILDARGS) -o=terrahelp-linux-amd64