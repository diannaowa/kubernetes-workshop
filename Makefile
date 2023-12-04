# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

# make all WHAT=cmd/serviceA
all: fmt vet  ## Build binary.
	hack/make-rules/build.sh $(WHAT)

# make run WHAT=cmd/serviceA
run: fmt vet ## Run a biz-splitting-tool from your host.
	go run $(WHAT)/main.go

fmt: ## Run go fmt against code.
	go fmt $(shell go list ./... | grep -v /vendor/)

vet: ## Run go vet against code.
	go vet $(shell go list ./... | grep -v /vendor/)