# ---------------------------------------------------------------------------- #
#             Apache 2.0 License Copyright © 2023 The Aurae Authors            #
#                                                                              #
#                +--------------------------------------------+                #
#                |   █████╗ ██╗   ██╗██████╗  █████╗ ███████╗ |                #
#                |  ██╔══██╗██║   ██║██╔══██╗██╔══██╗██╔════╝ |                #
#                |  ███████║██║   ██║██████╔╝███████║█████╗   |                #
#                |  ██╔══██║██║   ██║██╔══██╗██╔══██║██╔══╝   |                #
#                |  ██║  ██║╚██████╔╝██║  ██║██║  ██║███████╗ |                #
#                |  ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝ |                #
#                +--------------------------------------------+                #
#                                                                              #
#                         Distributed Systems Runtime                          #
#                                                                              #
# ---------------------------------------------------------------------------- #
#                                                                              #
#   Licensed under the Apache License, Version 2.0 (the "License");            #
#   you may not use this file except in compliance with the License.           #
#   You may obtain a copy of the License at                                    #
#                                                                              #
#       http://www.apache.org/licenses/LICENSE-2.0                             #
#                                                                              #
#   Unless required by applicable law or agreed to in writing, software        #
#   distributed under the License is distributed on an "AS IS" BASIS,          #
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.   #
#   See the License for the specific language governing permissions and        #
#   limitations under the License.                                             #
#                                                                              #
# ---------------------------------------------------------------------------- #

GORELEASER_FLAGS ?= --snapshot --rm-dist
all: compile

# Variables and Settings
version     ?=  0.0.1
target      ?=  ae
org         ?=  aurae-runtime
authorname  ?=  The Aurae Authors
authoremail ?=  info@aurae.io
license     ?=  Apache2
year        ?=  2023
copyright   ?=  Copyright (c) $(year)

COMMIT      := $(shell git rev-parse HEAD)
DATE        := $(shell date +%Y-%m-%d)
PKG_LDFLAGS := github.com/prometheus/common/version

compile: mod ## Compile for the local architecture ⚙
	@echo "Compiling..."
	go build -ldflags "\
	-X 'main.Version=$(version)' \
	-X 'main.AuthorName=$(authorname)' \
	-X 'main.AuthorEmail=$(authoremail)' \
	-X 'main.Copyright=$(copyright)' \
	-X 'main.License=$(license)' \
	-X 'main.Name=$(target)' \
	-X '${PKG_LDFLAGS}.Version=$(version)' \
	-X '${PKG_LDFLAGS}.BuildDate=$(DATE)' \
	-X '${PKG_LDFLAGS}.Revision=$(COMMIT)'" \
	-o bin/$(target) .

.PHONY: goreleaser
goreleaser: ## Run goreleaser directly at the pinned version 🛠
	go run github.com/goreleaser/goreleaser@v1.14 $(GORELEASER_FLAGS)

.PHONY: mod
mod: ## Go mod things
	go mod tidy
	go mod vendor
	go mod download

.PHONY: install
install: compile ## Install the program to /usr/bin 🎉
	@echo "Installing..."
	sudo cp bin/$(target) /usr/local/bin/$(target)

.PHONY: test
test: compile ## 🤓 Run go tests
	@echo "Testing..."
	go test -v ./...

.PHONY: clean
clean: ## Clean your artifacts 🧼
	@echo "Cleaning..."
	rm -rvf dist/*
	rm -rvf release/*

.PHONY: help
help:  ## Show help messages for make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: format
format: ## Format the code using gofmt
	@echo "Formatting..."
	@gofmt -s -w $(shell find . -name '*.go' -not -path "./vendor/*")

.PHONY: check-format
check-format: ## Used by CI to check if code is formatted
	@gofmt -l $(shell find . -name '*.go' -not -path "./vendor/*") | grep ".*" ; if [ $$? -eq 0 ]; then exit 1; fi

.PHONY: lint
lint: ## Runs the linter
	golangci-lint run

3PHONY: check-editorconfig
check-editorconfig: ## Use to check if the codebase follows editorconfig rules
	@docker run --rm --volume=$(shell PWD):/check mstruebing/editorconfig-checker
