.ONESHELL:

EXECUTABLES = git go find pwd
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))

ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

BINARY_NAME=waveman
PROJECT_NAME=waveman2
VERSION?=1.0.0
BUILD=`git rev-parse HEAD`
WINDOWS_PLATFORMS=windows/amd64 windows/386 windows/arm windows/arm64
UNIX_PLATFORMS=linux/amd64 linux/arm linux/arm64 linux/386 darwin/amd64
temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

ARTIFACTS_DIR=dist

# Setup linker flags option for build that interoperate with variable names in src code
LDFLAGS=-ldflags "-s -w -X 'main.Version=${VERSION}'"

release: $(UNIX_PLATFORMS) $(WINDOWS_PLATFORMS)

install-dev:
	CGO_ENABLED=0 go build $(LDFLAGS) -o /usr/local/bin/$(BINARY_NAME) main.go

clean:
	rm -rf dist/

.PHONY: clean release

$(UNIX_PLATFORMS): clean
	GOOS=$(os) GOARCH=$(arch) go build -v $(LDFLAGS) -o $(ARTIFACTS_DIR)/$(BINARY_NAME)_$(os)_$(arch)/$(BINARY_NANE)
	cd $(ARTIFACTS_DIR)/$(BINARY_NAME)_$(os)_$(arch)/
	tar -czf $(BINARY_NAME)_$(VERSION)_$(os)_$(arch).tar.gz $(BINARY_NANE)
	mv $(BINARY_NAME)_$(VERSION)_$(os)_$(arch).tar.gz ../
	cd ../
	rm -rf $(BINARY_NAME)_$(os)_$(arch)/

$(WINDOWS_PLATFORMS): clean
	GOOS=$(os) GOARCH=$(arch) go build -v $(LDFLAGS) -o $(ARTIFACTS_DIR)/$(BINARY_NAME)_$(os)_$(arch)/$(BINARY_NANE).exe
	cd $(ARTIFACTS_DIR)/$(BINARY_NAME)_$(os)_$(arch)/
	zip $(BINARY_NAME)_$(VERSION)_$(os)_$(arch).zip $(BINARY_NANE).exe
	mv $(BINARY_NAME)_$(VERSION)_$(os)_$(arch).zip ../
	cd ../
	rm -rf $(BINARY_NAME)_$(os)_$(arch)/

.PHONY: clean release