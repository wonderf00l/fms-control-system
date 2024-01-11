PROJECT_BIN = $(CURDIR)/bin
$(shell [ -f bin ] || mkdir -p $(PROJECT_BIN))
GOLANGCI_LINT = $(PROJECT_BIN)/golangci-lint

build:
	go build -o $(PROJECT_BIN)/app cmd/app/*.go

.PHONY:
run: build
	bin/app $(ARGS)

.install-linter:
	[ -f $(PROJECT_BIN)/golangci-lint ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(PROJECT_BIN) v1.55.2

.PHONY:
lint: .install-linter
	$(GOLANGCI_LINT) run ./... --config=configs/.golangci.yml

.PHONY:
lint-fast: .install-linter
	$(GOLANGCI_LINT) run ./... --fast --config=configs/.golangci.yml
