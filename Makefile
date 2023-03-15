SHELL := /bin/bash

NAME ?= terrahelp
BIN := $(CURDIR)/bin

go_files := $(shell find . -path '*/testdata' -prune -o -type f -name '*.go' -not -path "./vendor/*" -print)

.DEFAULT_GOAL := all
.PHONY := all vet fmt fmtcheck test install uninstall clean dependencies tidy-dependencies clean-dependencies

vet: $(go_files)
	go vet  ./...

fmt:
	@go run golang.org/x/tools/cmd/goimports -w $(go_files)

fmtcheck: $(go_files)
	# Checking format of Go files...
	@GOIMPORTS=$$(go run golang.org/x/tools/cmd/goimports -l $(go_files)) && \
	if [ "$$GOIMPORTS" != "" ]; then \
		go run golang.org/x/tools/cmd/goimports -d $(go_files); \
		exit 1; \
	fi

bin/.coverage.out: $(go_files)
	@mkdir -p bin/
	RS_API_URL=$(TEST_DB_URL) RS_USERNAME=$(TEST_USERNAME) RS_PASSWORD=$(TEST_PASSWORD) RS_DB=$(TEST_DB_NAME) go test -v ./... -coverpkg=$(shell go list ./... | xargs | sed -e 's/ /,/g') -coverprofile bin/.coverage.tmp
	@mv bin/.coverage.tmp bin/.coverage.out

test: bin/.coverage.out

coverage: bin/.coverage.out
	@go tool cover -html=bin/.coverage.out

bin/terrahelp: $(go_files)
	go build -trimpath -o ./bin/$(NAME)

build: vet fmtcheck bin/terrahelp

install: build
	cp -f bin/$(NAME) ${GOPATH}/bin/$(NAME)

uninstall:
	@ echo "==> Uninstalling $(NAME)"
	rm -f $$(which ${NAME})

clean:
	@ echo "==> Cleaning output files."
ifneq ($(BIN),)
	rm -rf $(BIN)
endif

dependencies:
	@ echo "==> Downloading dependencies for $(NAME)"
	@ go mod download

tidy-dependencies:
	@ echo "==> Tidying dependencies for $(NAME)"
	@ go mod tidy

clean-dependencies:
	@ echo "==> Cleaning dependencies for $(NAME)"
	@ rm -rf $(VENDOR)

all: test build
