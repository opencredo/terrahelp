test: deps
	go test -v ./...

build: deps
	go build -o ${GOPATH}/bin/terrahelp

