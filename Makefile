.ONESHELL:

EXECUTABLES = git go find pwd
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))

ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

BINARY=waveman2
VERSION=2.0.0
BUILD=`git rev-parse HEAD`
WINDOWS_PLATFORMS=windows/amd64 windows/386 windows/arm windows/arm64
UNIX_PLATFORMS=linux/amd64 linux/arm linux/arm64 linux/386 darwin/amd64
temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

ARTIFACTS_DIR=dist

# Setup linker flags option for build that interoperate with variable names in src code
# LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

release: $(UNIX_PLATFORMS) $(WINDOWS_PLATFORMS)

clean:
	rm -rf dist/

$(UNIX_PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) go build -v -o $(ARTIFACTS_DIR)/$(BINARY)_$(os)_$(arch)/waveman
	cd $(ARTIFACTS_DIR)/$(BINARY)_$(os)_$(arch)/
	tar -czf $(BINARY)_$(os)_$(arch).tar.gz waveman
	mv $(BINARY)_$(os)_$(arch).tar.gz ../
	rm -rf $(BINARY)_$(os)_$(arch)

$(WINDOWS_PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) go build -v -o $(ARTIFACTS_DIR)/$(BINARY)_$(os)_$(arch)/waveman.exe
	cd $(ARTIFACTS_DIR)/$(BINARY)_$(os)_$(arch)/
	zip $(BINARY)_$(os)_$(arch).zip waveman.exe
	mv $(BINARY)_$(os)_$(arch).zip ../
	rm -rf $(BINARY)_$(os)_$(arch)

.PHONY: clean release