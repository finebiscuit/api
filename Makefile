PKG := github.com/finebiscuit/api
BIN := bin/biscuit-api
VERSION ?= $(shell git describe --always --long --dirty)
BUILDTIME := $(shell date -u +'%Y-%m-%d_%T')
SOURCES := $(shell find . -name '*.go')
LDFLAGS :=

GO := go
GOBUILD := $(GO) build
GOGENERATE := $(GO) generate
GOTEST := $(GO) test
DOCKERCOMPOSE := docker-compose

$(BIN): $(SOURCES)
	@$(GOBUILD) -o $(BIN) \
		-trimpath \
		-ldflags="-X main.version=$(VERSION) -X main.buildTime=$(BUILDTIME) $(LDFLAGS)" \
		$(PKG)

.PHONY: generate
generate:
	@$(GOGENERATE) ./...

.PHONY: test
test:
	@$(GOTEST) ./...

.PHONY: postgres-up
postgres-up:
	@$(DOCKERCOMPOSE) up -d postgres

.PHONY: e2e
e2e: e2e-sqlite e2e-postgres

.PHONY: e2e-sqlite
e2e-sqlite:
	$(GOTEST) -tags sqlite $(PKG)/e2e

.PHONY: e2e-postgres
e2e-postgres: postgres-up
	@$(DOCKERCOMPOSE) run --rm api go test -tags postgres $(PKG)/e2e
