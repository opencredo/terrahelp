SHELL := /bin/bash

NAME ?= terrahelp

COMMIT = $(shell git rev-parse HEAD)

LDFLAGS := -ldflags "-X=main.commit=$(COMMIT)"
BUILDARGS ?= -mod=vendor

# Set to 1 to skip Go version check in the ensure-version target.
# This is useful when we want to build using development versions of Go.
SKIP_GO_REQ_VERSION_CHECK ?= 0

REQ_GO_VERSION := 1.13
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

.PHONY: ensure-version
ensure-version:
ifeq ($(SKIP_GO_REQ_VERSION_CHECK),1)
	@ echo "==> Skipping go version check"
else
	@ echo -n "==> Checking go version... "
ifeq ($(GO_VERSION),$(MAX_GO_VERSION))
	@ echo "OK!"
else
	@ $(error Found go $(GO_VERSION) but we require $(REQ_GO_VERSION))
endif
endif


.PHONY: check
check:
	go vet $(BUILDARGS) $(LDFLAGS) ./...

.PHONY: test
test:
	go test $(BUILDARGS) $(LDFLAGS) -v ./...

.PHONY: build
build: ensure-version check test
	go build $(BUILDARGS) $(LDFLAGS) -o bin/$(NAME)

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
$(PLATFORMS): ensure-version check test
	@ echo "==> Building $(OS) distribution"
	@ mkdir -p $(BIN)/$(OS)/$(ARCH)
	@ mkdir -p $(DIST)
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build $(BUILDARGS) $(LDFLAGS) -o $(BIN)/$(OS)/$(ARCH)/$(NAME)
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
