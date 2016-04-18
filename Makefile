.PHONY: build
.DEFAULT_GOAL := help

help: ## Displays this help message.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Generates a releasable binary.
	go build -ldflags="-s -w"
