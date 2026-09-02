GO ?= go
TOOLS_DIR = hack/tools

GOLANGCI_LINT_VER := 2.13.2
GOLANGCI_LINT := $(TOOLS_DIR)/golangci-lint-$(GOLANGCI_LINT_VER)

default: lint build test test-integration

.PHONY: build
build:
	$(GO) build ./...

.PHONY: test
test:
	$(GO) test ./cmd/... ./pkg/...

.PHONY: test-integration
test-integration:
	$(GO) test -count 1 ./test/integration/...

.PHONY: lint
lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run $(GOLANGCI_LINT_FLAGS) ./...

.PHONY: lint-fix
lint-fix: override GOLANGCI_LINT_FLAGS := $(GOLANGCI_LINT_FLAGS) --fix
lint-fix: lint

$(GOLANGCI_LINT):
	mkdir -p $(TOOLS_DIR)
	$(GO) tool codeberg.org/ntnn/mindl download -common -out $@ -tool golangci-lint -version $(GOLANGCI_LINT_VER)
	ln -sf $(notdir $@) $(TOOLS_DIR)/golangci-lint
