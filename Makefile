.ONESHELL:

EXECUTABLES = git go find pwd
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))

ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

BINARY=waveman2
VERSION?=1.0.0
BUILD=`git rev-parse HEAD`
WINDOWS_PLATFORMS=windows/amd64 windows/386 windows/arm windows/arm64
UNIX_PLATFORMS=linux/amd64 linux/arm linux/arm64 linux/386 darwin/amd64
temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

ARTIFACTS_DIR=dist

# Setup linker flags option for build that interoperate with variable names in src code
LDFLAGS=-ldflags "-s -w -X main.Version=${VERSION}"

release: $(UNIX_PLATFORMS) $(WINDOWS_PLATFORMS)

clean:
	rm -rf dist/

.PHONY: clean release

$(UNIX_PLATFORMS): clean
	GOOS=$(os) GOARCH=$(arch) go build -v $(LDFLAGS) -o $(ARTIFACTS_DIR)/$(BINARY)_$(os)_$(arch)/assignmentctl
	cd $(ARTIFACTS_DIR)/$(BINARY)_$(os)_$(arch)/
	tar -czf $(BINARY)_$(VERSION)_$(os)_$(arch).tar.gz assignmentctl
	mv $(BINARY)_$(VERSION)_$(os)_$(arch).tar.gz ../
	cd ../
	rm -rf $(BINARY)_$(os)_$(arch)/

$(WINDOWS_PLATFORMS): clean
	GOOS=$(os) GOARCH=$(arch) go build -v $(LDFLAGS) -o $(ARTIFACTS_DIR)/$(BINARY)_$(os)_$(arch)/assignmentctl.exe
	cd $(ARTIFACTS_DIR)/$(BINARY)_$(os)_$(arch)/
	zip $(BINARY)_$(VERSION)_$(os)_$(arch).zip assignmentctl.exe
	mv $(BINARY)_$(VERSION)_$(os)_$(arch).zip ../
	cd ../
	rm -rf $(BINARY)_$(os)_$(arch)/


.PHONY: clean release