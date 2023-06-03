BINARY_NAME := ssmbrowse
CURRENT_DATETIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION ?= "dev"

build:
	@CGO_ENABLED=1 go build  \
		-ldflags "-X main.version=$(VERSION) -X main.date=$(CURRENT_DATETIME)" \
		-o $(BINARY_NAME) *.go
	

run: build
	@./$(BINARY_NAME)

lint:
	@gofmt -s -w .

release:
	@goreleaser release --clean

release-snapshot:
	@goreleaser release --snapshot --clean
