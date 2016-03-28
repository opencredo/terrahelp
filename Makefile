deps:
	go get -u github.com/kardianos/govendor
	govendor sync

test: deps
	go test

build: deps
	go build -o ${GOPATH}/bin/terrahelp

