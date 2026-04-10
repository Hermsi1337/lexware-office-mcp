.PHONY: build release-check release-snapshot

GORELEASER_VERSION := v2.15.2
GORELEASER_IMAGE := goreleaser/goreleaser:$(GORELEASER_VERSION)
GORELEASER_CONFIG := build/goreleaser/.goreleaser.yml

DOCKER_RUN := docker run --rm \
	-v $(CURDIR):/workspace \
	-v /var/run/docker.sock:/var/run/docker.sock \
	-w /workspace \
	$(GORELEASER_IMAGE)

VERSION ?= dev

build:
	go build -ldflags "-X github.com/dennis/lexware-office-mcp/internal/version.Version=$(VERSION)" -o bin/lexware-office-mcp ./cmd/lexware-office-mcp

release-check:
	$(DOCKER_RUN) check --config $(GORELEASER_CONFIG)

release-snapshot:
	$(DOCKER_RUN) release --snapshot --clean --config $(GORELEASER_CONFIG)
