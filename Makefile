BINARY_NAME := monitor
PKG := github.com/yektasrk/http-monitor
GO ?= go

.PHONY: all dep build clean

all: build

dep: ## Get the dependencies
	$(GO) mod download

build: ## Build the binary file
	$(GO) build -v $(PKG)/cmd/$(BINARY_NAME)d

clean: ## Remove previous build
	@rm -f $(BINARY_NAME)d

help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
