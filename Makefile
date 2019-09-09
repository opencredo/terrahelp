SHELL := /bin/bash

NAME ?= terrahelp

BUILDARGS ?= -mod=vendor

BIN := $(CURDIR)/bin
DIST := $(CURDIR)/dist
OUPUT_FILES := $(BIN) $(DIST)

PLATFORMS ?= darwin linux
ARCH ?= amd64
OS = $(word 1, $@)

check:
	go vet $(BUILDARGS) ./...
.PHONY: check

test:
	go test $(BUILDARGS) -v ./...
.PHONY: test

build:
	go build $(BUILDARGS) -o bin/$(NAME)
.PHONY: build

install:
	go build $(BUILDARGS) -o ${GOPATH}/bin/$(NAME)
.PHONY: install

uninstall:
	@ echo "==> Uninstalling $(NAME)"
	rm -f $$(which ${NAME})
.PHONY: uninstall

clean:
	@ echo "==> Cleaning output files."
ifneq ($(OUPUT_FILES),)
	rm -rf $(OUPUT_FILES)
endif
.PHONY: clean

$(PLATFORMS):
	@ echo "==> Building $(OS) distribution"
	@ mkdir -p $(BIN)/$(OS)/$(ARCH)
	@ mkdir -p $(DIST)
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build $(BUILDARGS) -o $(BIN)/$(OS)/$(ARCH)/$(NAME)
	cp -f $(BIN)/$(OS)/$(ARCH)/$(NAME) $(DIST)/$(NAME)-$(OS)-$(ARCH)
.PHONY: $(PLATFORMS)

dist: $(PLATFORMS)
	@ true
.PHONY: dist

dependencies:
	@ echo "==> Downloading dependencies for $(TARGET)"
	@ go mod download
.PHONY: dependencies

vendor-dependencies:
	@ echo "==> Downloading dependencies for $(TARGET)"
	@ go mod vendor
.PHONY: vendor-dependencies

tidy-dependencies:
	@ echo "==> Tidying dependencies for $(TARGET)"
	@ go mod tidy
.PHONY: tidy-dependencies

clean-dependencies:
	@ echo "==> Cleaning dependencies for $(TARGET)"
	@ rm -rf $(VENDOR)
.PHONY: clean-dependencies