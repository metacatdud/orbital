
# Version Vars
VERSION_TAG := $(shell git describe --tags --always)
VERSION_VERSION := $(shell git log --date=iso --pretty=format:"%cd" -1) $(VERSION_TAG)
VERSION_COMPILE := $(shell date +"%F %T %z") by $(shell go version)
VERSION_BRANCH  := $(shell git rev-parse --abbrev-ref HEAD)
VERSION_GIT_DIRTY := $(shell git diff --no-ext-diff 2>/dev/null | wc -l | awk '{print $1}')
VERSION_DEV_PATH:= $(shell pwd)

# Go Checkup
GOPATH ?= $(shell go env GOPATH)
GO111MODULE:=auto
export GO111MODULE
ifeq "$(GOPATH)" ""
  $(error Please set the environment variable GOPATH before running `make`)
endif
PATH := ${GOPATH}/bin:$(PATH)
GCFLAGS=-gcflags="all=-trimpath=${GOPATH}"
LDFLAGS=-ldflags="-s -w -X 'main.Version=${VERSION_TAG}' -X 'main.Compile=${VERSION_COMPILE}' -X 'main.Branch=${VERSION_BRANCH}'"

GO = go

V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m➡\033[0m")

# Commands
.PHONY: all
all: | init deps

.PHONY: init
init: ; $(info $(M) Installing tools dependencies ...) @ ## Install tools dependencies
	@if ! pre-commit --version > /dev/null 2>&1; then \
		echo "pre-commit is not installed. Please install it using one of the following methods:"; \
		echo "- For Debian/Ubuntu-based systems: apt install pre-commit"; \
		echo "- For macOS (Homebrew): brew install pre-commit"; \
		echo "- For Python environments: pip install pre-commit"; \
		exit 1; \
	fi

	pre-commit install --install-hooks
	pre-commit install --hook-type commit-msg

	$Q $(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$Q $(GO) install github.com/go-critic/go-critic/cmd/gocritic@latest
	$Q $(GO) install github.com/sqs/goreturns@latest

.PHONY: deps
deps: ; $(info $(M) Installing project dependencies ...) @ ## Install project dependencies
	$Q $(GO) mod tidy

.PHONY: test
test:  ; $(info $(M) Running unit tests ...)	@ ## Run unit tests
	$Q $(GO) test -v  -coverprofile coverage.out ./...

.PHONY: build
build: ; $(info $(M) Building executable...) @ ## Build program binary
	$Q echo "ver   : ${VERSION_TAG}"
	$Q echo "veriso: ${VERSION_VERSION}"
	$Q echo "vergo : ${VERSION_COMPILE}"
	$Q mkdir -p bin
	$Q $(GO) generate ./...
	$Q ret=0 && for d in $$($(GO) list -f '{{if (eq .Name "main")}}{{.ImportPath}}{{end}}' ./...); do \
		b=$$(basename $${d}) ; \
		$(GO) build ${LDFLAGS} ${GCFLAGS} -o bin/$${b} $$d || ret=$$? ; \
		echo "$(M) Build: bin/$${b}" ; \
		echo "$(M) Done!" ; \
	done ; exit $$ret

.PHONY: release
release: ; $(info $(M) Buildig release version) @ ## Build OS specific release version
	$Q rm -rf release
	$Q mkdir -p release
	$Q ret=0 && for os in darwin linux windows; do \
		for arch in amd64 arm64; do \
			binary_name="atomika_$${os}_$${arch}_${VERSION_TAG}" ; \
			ext="" ; \
			if [ "$$os" = "windows" ]; then ext=".exe"; fi ; \
			GOOS=$$os GOARCH=$$arch $(GO) build ${LDFLAGS} ${GCFLAGS} -o release/$$binary_name$$ext || ret=$$? ; \
			echo "$(M) Build: release/$$binary_name$$ext" ; \
		done; \
		if [ "$$os" = "windows" ]; then \
			if which zip >/dev/null 2>&1; then \
				echo "$(M) Archiving release files for $$os using zip..." ; \
				zip -j release/atomika_$${os}_${VERSION_TAG}.zip release/atomika_$${os}_*_$${VERSION_TAG}* ; \
				echo "$(M) Archive created: release/atomika_$${os}_${VERSION_TAG}.zip" ; \
			else \
				echo "$(M) Zip command not found, skipping archive for windows." ; \
			fi; \
		else \
			echo "$(M) Archiving release files for $$os using tar..." ; \
			find release -name "atomika_$${os}_*_${VERSION_TAG}*" -print | xargs tar -czf release/atomika_$${os}_${VERSION_TAG}.tar.gz || ret=$$? ; \
			echo "$(M) Archive created: release/atomika_$${os}_${VERSION_TAG}.tar.gz" ; \
		fi; \
		echo "$(M) Cleaning up binary files for $$os..." ; \
		find release -name "atomika_$${os}_*_${VERSION_TAG}*" -type f -not -name "*.zip" -not -name "*.tar.gz" -exec rm {} \; ; \
	done; exit $$ret

.PHONY: run
run: ; $(info $(M) Running dev build (on the fly) ...) @ ## Run intermediate builds
	$Q $(GO) run -race ./...

help:
	$Q echo "\nAtomika.io\n----------------"
	$Q grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
