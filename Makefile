
# planned to use docker also... should be rootless, and podman is my prefered OCI engine
OCI=podman
OCI_OPTS=--userns=keep-id
GO_IMAGE=docker.io/golang:1.22
GO_BUILD_OPTS=-e CGO_ENABLED=0

# Where to install the binary, default to tht HOME/.local/bin as it is the recommended path for user binaries
# by freedesktop.org. Should be OK for most Linux distributions (even with WSL2), FreeBSD and OpenBSD.
# I'm not sure how to handle the case of MacOS or Windows.
PREFIX=~/.local

VERSION:=$(shell git describe --tags --always --dirty)
# please, do not change the BUILD variable, it is used to identify the build type
BUILD=local-
GO_VERSION_LD:=-X main.Version=$(BUILD)$(VERSION)


all: dist/pifpaf.linux.amd64 \
	dist/pifpaf.windows.amd64 \
	dist/pifpaf.darwin.amd64 \
	dist/pifpaf.freebsd.amd64 \
	dist/pifpaf.linux.arm64 \
	dist/pifpaf.freebsd.arm64

dist/pifpaf%:
	# split the target into the name and the extension
	$(eval TARGET := $(subst ., ,$@))
	$(eval NAME := $(word 1, $(TARGET)))
	$(eval GOOS := $(word 2, $(TARGET)))
	$(eval GOARCH := $(word 3, $(TARGET)))
	echo "building $(NAME) for $(GOOS) : $(GOARCH)"
	$(OCI) run --rm -v $(PWD):/src:z -w /src \
		-e GOOG=$(GOOS) \
		-e GOARCH=$(GOARCH) \
		$(GO_BUILD_OPTS) \
		$(OCI_OPTS) \
		$(GO_IMAGE) \
		go build -ldflags="$(GO_VERSION_LD)"  -o $@ ./cmd/pifpaf
	strip $@ || :

install:
	@ARCH=$(shell uname -m); \
	[ $$ARCH == "x86_64" ] && ARCH="amd64"; \
	OS=$(shell uname -s); \
	[ $$OS == "Darwin" ] && OS="darwin"; \
	[ $$OS == "FreeBSD" ] && OS="freebsd"; \
	[ $$OS == "Linux" ] && OS="linux"; \
	[ $$OS == "Windows" ] && OS="windows"; \
	$(MAKE) dist/pifpaf.$$OS.$$ARCH;\
	install -Dm755 dist/pifpaf.$$OS.$$ARCH $(PREFIX)/bin/pifpaf


clean:
	rm -rf dist/*
