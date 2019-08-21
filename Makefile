BUILDARGS ?= -mod=vendor

test:
	go test $(BUILDARGS) -v ./...

build:
	go build $(BUILDARGS) -o ${GOPATH}/bin/terrahelp

