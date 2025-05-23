SHELL = bash

export PWD := $(shell pwd)

export ARCH := $(shell uname -m)
export ARCH_DIST := $(ARCH)
ifeq ($(ARCH), x86_64)
	ARCH_DISK := amd64
endif

export OS_NAME = $(shell uname -s)
export OS_TYPE := linux
ifeq ($(OS_NAME), Darwin)
	OS_TYPE := darwin
endif

export BUILD_AT := $(shell date -u +'%Y-%m-%dT%T%Z')

# Tag of the current commit, if any. If this is not "" then we are building a release
export RELEASE_TAG := $(shell git tag -l --points-at HEAD| sort | head -n 1)
# Last tag on this branch
export LAST_TAG := $(shell git describe --tags --abbrev=0)
export BUILD_VERSION := $(or $(RELEASE_TAG), $(LAST_TAG))
export TAG_BRANCH := .$(BRANCH)
# If building HEAD or main then unset TAG_BRANCH
ifeq ($(subst HEAD,,$(subst main,,$(BRANCH))),)
	TAG_BRANCH :=
endif

# COMMIT is the commit hash
export COMMIT := $(shell git log -1 --format="%H" | head -1)
# COMMIT_NUMBER is the number commits since last tag.
export COMMIT_NUMBER := $(shell git rev-list --count $(RELEASE_TAG)...HEAD)

# Make version suffix -NNNN.CCCCCCCC (N=Commit number, C=Commit)
export VERSION_SUFFIX := $(COMMIT_NUMBER).$(shell git show --no-patch --no-notes --pretty='%h' HEAD)
export VERSION := $(RELEASE_TAG)-$(VERSION_SUFFIX)$(TAG_BRANCH)

# Pass in GOTAGS=xyz on the make command line to set build tags
ifdef GOTAGS
	BUILDTAGS=-tags "$(GOTAGS)"
	LINTTAGS=--build-tags "$(GOTAGS)"
endif

export OUTPUT_DIR := $(PWD)/_output
export BIN := $(OUTPUT_DIR)/bin

buildVersionLDFlag := -X github.com/diabloneo/s3-mcp/pkg/common.Version=$(BUILD_VERSION) \
	-X github.com/diabloneo/s3-mcp/pkg/common.GitSha=$(COMMIT) \
	-X github.com/diabloneo/s3-mcp/pkg/common.BuildTime=$(BUILD_AT)

.PHONY: build
build:
	CGO_ENABLED=0 go build -trimpath -ldflags "$(buildVersionLDFlag)" $(BUILDTAGS) -o $(BIN)/s3-mcp github.com/diabloneo/s3-mcp/cmd/s3-mcp

.PHONY: test
test:
	go test -timeout 1h `go list ./...`

.PHONY: validate
validate: 
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.6
	@echo "Lint code with golangci-lint"
	PATH=$(shell go env GOPATH)/bin:$(PATH) golangci-lint run
