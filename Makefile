BINDIR := ./bin
SRCDIR := ./cmd
API_TARGET := $(BINDIR)/api
WEB_TARGET := $(BINDIR)/web

GO := go
GOFLAGS :=

.PHONY: all
all: build-api build-web

.PHONY: build-api
build-api: $(API_TARGET)

$(API_TARGET): $(SRCDIR)/gowiki_api/main.go
	@echo "Building API server..."
	$(GO) build $(GOFLAGS) -o $(API_TARGET) $(SRCDIR)/gowiki_api

.PHONY: build-web
build-web: $(WEB_TARGET)

$(WEB_TARGET): $(SRCDIR)/gowiki_web/main.go
	@echo "Cleaning up..."
	rm -rf $(WEB_TARGET)
	@echo "Building Web server..."
	$(GO) build $(GOFLAGS) -o $(WEB_TARGET) $(SRCDIR)/gowiki_web

.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf $(BINDIR)

.PHONY: help
help:
	@echo "Usage:"
	@echo "  make build-api      Build the API server and place it in $(BINDIR)"
	@echo "  make build-web      Build the Web server and place it in $(BINDIR)"
	@echo "  make clean          Remove all build artifacts"
