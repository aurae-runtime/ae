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

all: test

.PHONY: proto
proto: ## Generate code from protobuf schemas
	@buf --version >/dev/null 2>&1 || (echo "Warning: buf is not installed! Please install the 'buf' command line tool: https://docs.buf.build/installation"; exit 1)
	buf generate -v https://github.com/aurae-runtime/aurae.git#branch=main,subdir=api

.PHONY: mod
mod: ## Go mod things
	go mod tidy
	go mod vendor
	go mod download

.PHONY: test
test: clean proto ## 🤓 Run go tests
	@echo "Testing..."
	go test -v ./...

.PHONY: fmt
fmt:
	@gofumpt -version >/dev/null 2>&1 || (echo "Warning: gofumpt is not installed! Please install the 'gofumpt' command line tool: https://pkg.go.dev/mvdan.cc/gofumpt"; exit 1)
	gofumpt -l -w .

.PHONY: lint
lint: fmt ## Linting the code according to the styling guide.
	@golangci-lint --version >/dev/null 2>&1 || (echo "Warning: golangci-lint is not installed! Please install the 'golangci-lint' command line tool: https://golangci-lint.run/usage/install/"; exit 1)
	golangci-lint run -c .golangci.yaml

.PHONY: clean
clean: ## Clean your artifacts 🧼
	@echo "Cleaning..."
	rm -rvf pkg/api/

.PHONY: help
help:  ## Show help messages for make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'
