DIST_DIR ?= dist

.PHONY: all
all: clean unit_test build

#
#	Release targets
#

.PHONY: dry_release
dry_release: clean
	goreleaser --skip-publish

.PHONY: install
install:
	mv $(DIST_DIR)/sind $(GOPATH)/bin/sind

#
# Test targets
#

.PHONY: unit_test
unit_test:
	go test -v -race -timeout=30s -cover -run=$(T) ./pkg/...

.PHONY: lint
lint:
	golangci-lint run

#
# Build targets
#

.PHONY: download
download:
	go mod download

.PHONY: build
build: clean $(DIST_DIR) binary

.PHONY: binary
binary:
	CGO_ENABLED=0 go build -ldflags='-s -w' -o $(DIST_DIR)/goats ./cmd/goats

$(DIST_DIR):
	mkdir -p $(DIST_DIR)

.PHONY: clean
clean:
	rm -rf $(DIST_DIR)
