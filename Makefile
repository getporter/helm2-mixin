MIXIN = helm
PKG = github.com/deislabs/porter-$(MIXIN)

COMMIT ?= $(shell git rev-parse --short HEAD)
VERSION ?= $(shell git describe --tags --dirty='+dev' --abbrev=0 2> /dev/null || echo v0)
PERMALINK ?= $(shell git name-rev --name-only --tags --no-undefined HEAD &> /dev/null && echo latest || echo canary)

LDFLAGS = -w -X $(PKG)/pkg.Version=$(VERSION) -X $(PKG)/pkg.Commit=$(COMMIT)
XBUILD = CGO_ENABLED=0 go build -a -tags netgo -ldflags '$(LDFLAGS)'
BINDIR = bin/mixins/$(MIXIN)

CLIENT_PLATFORM = $(shell go env GOOS)
CLIENT_ARCH = $(shell go env GOARCH)
RUNTIME_PLATFORM = linux
RUNTIME_ARCH = amd64
SUPPORTED_CLIENT_PLATFORMS = linux darwin windows
SUPPORTED_CLIENT_ARCHES = amd64 386

ifeq ($(CLIENT_PLATFORM),windows)
FILE_EXT=.exe
else ifeq ($(RUNTIME_PLATFORM),windows)
FILE_EXT=.exe
else
FILE_EXT=
endif

REGISTRY ?= $(USER)

build: build-client build-runtime

build-runtime:
	mkdir -p $(BINDIR)
	GOARCH=$(RUNTIME_ARCH) GOOS=$(RUNTIME_PLATFORM) $(XBUILD) -o $(BINDIR)/$(MIXIN)-runtime$(FILE_EXT) ./cmd/$(MIXIN)

build-client:
	mkdir -p $(BINDIR)
	go build -o $(BINDIR)/$(MIXIN)$(FILE_EXT) ./cmd/$(MIXIN)

build-all: build-runtime $(addprefix build-for-,$(SUPPORTED_CLIENT_PLATFORMS))
	cp $(BINDIR)/$(MIXIN)-runtime$(FILE_EXT) $(BINDIR)/$(VERSION)

build-for-%:
	$(MAKE) CLIENT_PLATFORM=$* xbuild-client

xbuild-client: $(BINDIR)/$(VERSION)/$(MIXIN)-$(CLIENT_PLATFORM)-$(CLIENT_ARCH)$(FILE_EXT)
$(BINDIR)/$(VERSION)/$(MIXIN)-$(CLIENT_PLATFORM)-$(CLIENT_ARCH)$(FILE_EXT):
	mkdir -p $(dir $@)
	$(XBUILD) -o $@ ./cmd/$(MIXIN)

test: test-unit

test-unit: build
	go test ./...

publish:
	cp -R $(BINDIR)/$(VERSION) $(BINDIR)/$(PERMALINK)
	# AZURE_STORAGE_CONNECTION_STRING will be used for auth in the following command
	az storage blob upload-batch -d porter/mixins/$(MIXIN)/$(VERSION) -s $(BINDIR)/$(VERSION)

clean:
	-rm -fr bin/
