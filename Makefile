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

compile: mod ## Compile for the local architecture ⚙
	@echo "Compiling..."
	go build -ldflags "\
	-X 'main.Version=$(version)' \
	-X 'main.AuthorName=$(authorname)' \
	-X 'main.AuthorEmail=$(authoremail)' \
	-X 'main.Copyright=$(copyright)' \
	-X 'main.License=$(license)' \
	-X 'main.Name=$(target)'" \
	-o bin/$(target) .

mod: ## Go mod things
	go mod tidy
	go mod vendor
	go mod download

install: ## Install the program to /usr/bin 🎉
	@echo "Installing..."
	sudo cp $(target) /usr/bin/$(target)

test: clean compile install ## 🤓 Run go tests
	@echo "Testing..."
	go test -v ./...

clean: ## Clean your artifacts 🧼
	@echo "Cleaning..."
	rm -rvf release/*

.PHONY: release
release: ## Make the binaries for a GitHub release 📦
	mkdir -p release
	GOOS="linux" GOARCH="amd64" go build -ldflags "-X 'main.Version=$(version)'" -o release/$(target)-linux-amd64 cmd/*.go
	GOOS="linux" GOARCH="arm" go build -ldflags "-X 'main.Version=$(version)'" -o release/$(target)-linux-arm cmd/*.go
	GOOS="linux" GOARCH="arm64" go build -ldflags "-X 'main.Version=$(version)'" -o release/$(target)-linux-arm64 cmd/*.go
	GOOS="linux" GOARCH="386" go build -ldflags "-X 'main.Version=$(version)'" -o release/$(target)-linux-386 cmd/*.go
	GOOS="darwin" GOARCH="amd64" go build -ldflags "-X 'main.Version=$(version)'" -o release/$(target)-darwin-amd64 cmd/*.go

.PHONY: help
help:  ## Show help messages for make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'

