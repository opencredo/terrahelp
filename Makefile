SHELL := /bin/bash

NAME ?= terrahelp

BUILDARGS ?= -mod=vendor

REQ_GO_VERSION := 1.16
GO_VERSION := $(shell go version | sed -E 's/^go version go([0-9]+.[0-9]+.[0-9]+).*$$/\1/')
MAX_GO_VERSION := $(shell printf "%s\n%s" $(REQ_GO_VERSION) $(GO_VERSION) | sort -V -r | head -1)

BIN := $(CURDIR)/bin
DIST := $(CURDIR)/dist
OUPUT_FILES := $(BIN) $(DIST)

PLATFORMS ?= darwin linux
ARCH ?= amd64
OS = $(word 1, $@)

SHA256_CMD = sha256sum
ifeq ($(shell uname), Darwin)
	SHA256_CMD = shasum -a 256
endif

.PHONY: check
check:
	go vet $(BUILDARGS) ./...

.PHONY: test
test:
	go test $(BUILDARGS) -v ./...

.PHONY: build
build: check test
	go build $(BUILDARGS) -o bin/$(NAME)

.PHONY: install
install: build
	cp -f bin/$(NAME) ${GOPATH}/bin/$(NAME)

.PHONY: uninstall
uninstall:
	@ echo "==> Uninstalling $(NAME)"
	rm -f $$(which ${NAME})

.PHONY: clean
clean:
	@ echo "==> Cleaning output files."
ifneq ($(OUPUT_FILES),)
	rm -rf $(OUPUT_FILES)
endif

.PHONY: $(PLATFORMS)
$(PLATFORMS): check test
	@ echo "==> Building $(OS) distribution"
	@ mkdir -p $(BIN)/$(OS)/$(ARCH)
	@ mkdir -p $(DIST)

	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build $(BUILDARGS) -o $(BIN)/$(OS)/$(ARCH)/$(NAME)
	cp -f $(BIN)/$(OS)/$(ARCH)/$(NAME) $(DIST)/$(NAME)-$(OS)-$(ARCH)

	@ $(SHA256_CMD) $(DIST)/$(NAME)-$(OS)-$(ARCH) | awk '{$$2=" $(NAME)-$(OS)-$(ARCH)"; print $$0}' >> $(DIST)/$(NAME).SHA256SUMS

.PHONY: dist
dist: $(PLATFORMS)
	@ touch $(DIST)/$(NAME).SHA256SUMS

.PHONY: dependencies
dependencies:
	@ echo "==> Downloading dependencies for $(NAME)"
	@ go mod download

.PHONY: vendor-dependencies
vendor-dependencies:
	@ echo "==> Downloading dependencies for $(NAME)"
	@ go mod vendor

.PHONY: tidy-dependencies
tidy-dependencies:
	@ echo "==> Tidying dependencies for $(NAME)"
	@ go mod tidy

.PHONY: clean-dependencies
clean-dependencies:
	@ echo "==> Cleaning dependencies for $(NAME)"
	@ rm -rf $(VENDOR)
