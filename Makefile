TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
COVER_TEST?=$$(go list ./... |grep -v 'vendor')

PKG_OS ?= darwin linux
PKG_ARCH ?= amd64
BASE_PATH ?= $(shell pwd)
BUILD_PATH ?= $(BASE_PATH)/build
PROVIDER := $(shell basename $(BASE_PATH))
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
VERSION ?= v0.0.0
ifneq ($(origin TRAVIS_TAG), undefined)
	BRANCH := $(TRAVIS_TAG)
	VERSION := $(TRAVIS_TAG)
endif

SYSTEM_ARCH ?= amd64
SYSTEM_OS ?= linux
ifeq ($(OS),Windows_NT)
    SYSTEM_OS := windows
    ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
        SYSTEM_ARCH := amd64
    else
        ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
            SYSTEM_ARCH := amd64
        endif
        ifeq ($(PROCESSOR_ARCHITECTURE),x86)
            SYSTEM_ARCH := i386
        endif
    endif
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        SYSTEM_OS := linux
    endif
    ifeq ($(UNAME_S),Darwin)
        SYSTEM_OS := darwin
    endif
    UNAME_M := $(shell uname -m)
    ifeq ($(UNAME_M),x86_64)
        SYSTEM_ARCH := amd64
    endif
    ifneq ($(filter %86,$(UNAME_M)),)
        SYSTEM_ARCH := i386
    endif
    ifneq ($(filter arm%,$(UNAME_M)),)
        SYSTEM_ARCH := arm
    endif
endif

default: build

build: fmtcheck
	go build -v .

local-install: build
	mkdir -p ~/.terraform.d/plugins/$(SYSTEM_OS)_$(SYSTEM_ARCH)/; \
	mv ./terraform-provider-online ~/.terraform.d/plugins/$(SYSTEM_OS)_$(SYSTEM_ARCH)/

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

testrace: fmtcheck
	TF_ACC= go test -race $(TEST) $(TESTARGS)

cover:
	@go tool cover 2>/dev/null; if [ $$? -eq 3 ]; then \
		go get -u golang.org/x/tools/cmd/cover; \
	fi
	go test $(COVER_TEST) -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

test-compile: fmtcheck
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./terraform-provider-online"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

packages:
	@for os in $(PKG_OS); do \
		for arch in $(PKG_ARCH); do \
			mkdir -p $(BUILD_PATH)/$(PROVIDER)_$${os}_$${arch} && \
			cd $(BASE_PATH) && \
			cgo_enabled=0 GOOS=$${os} GOARCH=$${arch} go build -o $(BUILD_PATH)/$(PROVIDER)_$${os}_$${arch}/$(PROVIDER)_$(VERSION) . && \
			cd $(BUILD_PATH) && \
			tar -cvzf $(BUILD_PATH)/$(PROVIDER)_$(BRANCH)_$${os}_$${arch}.tar.gz $(PROVIDER)_$${os}_$${arch}/; \
		done; \
	done;

clean:
	@rm -rf $(BUILD_PATH)

.PHONY: build test testacc testrace cover vet fmt fmtcheck errcheck test-compile
