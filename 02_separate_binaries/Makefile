SHELL=/bin/bash -o pipefail

GO_VERSION ?= 1.15
NAMESPACE=sbogacz
APPNAME=going-serverless
TALK_PHASE=02_separate_binaries
BINARY ?= handler

all: build

.PHONY: build
build: build-lambda

.PHONY: build-local
build-local:
	go build -o $(BINARY) \
		github.com/$(NAMESPACE)/$(APPNAME)/$(TALK_PHASE)/cmd/http

.PHONY: build-lambda
build-lambda:
	GOOS=linux go build -o $(BINARY) \
		github.com/$(NAMESPACE)/$(APPNAME)/$(TALK_PHASE)/cmd/lambda
	zip $(BINARY).zip $(BINARY)

.PHONY: clean
clean: 
	go clean
	rm -rf $(BINARY) $(BINARY).zip
