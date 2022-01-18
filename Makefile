.PHONY: build check test

GOOS ?= linux
GOARCH ?= amd64
VERSION ?= LOCAL

pkgs := $(shell go list ./...)
testPkgs := $(shell go list ./... | grep -v /test/e2e)

build: generate-mocks test
	@echo "make build"
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o build/cardano-tool cmd/cardano-tool.go

check: fmt vet lint

fmt:
	@echo "running go fmt"
	go fmt ./...

vet:
	@echo "running go vet"
	go vet $(pkgs)

lint:
	revive -config revive_conf.toml $(pkgs)

test: check
	go test $(testPkgs)

func-test: build
	@echo "running ginkgo"
	ginkgo -r test

send:
	@echo "make send"
	cd cmd/send-tool && CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ../../build/send-tool

generate-mocks:
	mockgen -source=pkg/cli/root_cmd.go -destination=pkg/mocks/cli/root_cmd.go

dependencies: check-system-dependencies
ifeq (, $(shell which revive))
	@echo "== cannot find revive installing"
	go install github.com/mgechev/revive@latest
endif

check-system-dependencies:
	@echo "== check-system-dependencies"
ifeq (, $(shell which go))
	$(error "golang not found in PATH")
endif