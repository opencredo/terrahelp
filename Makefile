SHELL := /bin/bash

TARGET ?= $(shell echo $${PWD\#\#*/})

VENDOR := $(CURDIR)/vendor

BIN := $(CURDIR)/bin
DIST := $(CURDIR)/dist
OUPUT_FILES := $(BIN) $(DIST)

PLATFORMS ?= darwin linux
ARCH ?= amd64
OS = $(word 1, $@)

VERSION ?= vlocal
COMMIT = $(shell git rev-parse HEAD)

LDFLAGS := -ldflags "-X=main.version=$(VERSION)"

# Go source files, excluding vendor directory
SRC := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

$(TARGET): $(SRC)
	@ mkdir -p $(BIN)
	@ go build $(LDFLAGS) -o $(BIN)/$(TARGET) .

build: $(TARGET)
	@ echo "==> Building $(TARGET)"
.PHONY: build

test: dependencies
	@ echo "==> Testing $(TARGET)"
	@ go test -v ./...
.PHONY: test

install:
	@ echo "==> Installing $(TARGET)"
	@ go install $(LDFLAGS)
.PHONY: install

uninstall: clean
	@ echo "==> Uninstalling $(TARGET)"
	rm -f $$(which ${TARGET})
.PHONY: uninstall

$(PLATFORMS):
	@ echo "==> Building $(OS) distribution"
	@ mkdir -p $(BIN)/$(OS)/$(ARCH)
	@ mkdir -p $(DIST)
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build $(LDFLAGS) -o $(BIN)/$(OS)/$(ARCH)/$(TARGET)
	cp -f $(BIN)/$(OS)/$(ARCH)/$(TARGET) $(DIST)/$(TARGET)-$(OS)-$(ARCH)
.PHONY: $(PLATFORMS)

dist: $(PLATFORMS)
	@ true
.PHONY: dist

check:
	@ echo "==> Checking $(TARGET)"
	@ go vet ./...
.PHONY: check

clean:
	@ echo "==> Cleaning output files."
ifneq ($(OUPUT_FILES),)
	rm -rf $(OUPUT_FILES)
endif
.PHONY: clean

dependencies:
	@ echo "==> Downloading dependencies for $(TARGET)"
	@ go mod download
.PHONY: dependencies

vendor-dependencies:
	@ echo "==> Downloading dependencies for $(TARGET)"
	@ go mod vendor
.PHONY: vendor-dependencies

clean-dependencies:
	@ echo "==> Cleaning dependencies for $(TARGET)"
	@ rm -rf $(VENDOR)
.PHONY: clean-dependencies