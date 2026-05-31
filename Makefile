.PHONY: check fmt test test-install vet diff-check

export GOCACHE := $(CURDIR)/.cache/go-build
export GOMODCACHE := $(CURDIR)/.cache/go-mod

check: fmt test test-install vet diff-check

fmt:
	mkdir -p $(GOCACHE) $(GOMODCACHE)
	gofmt -w ./cmd ./internal

test:
	go test ./...

test-install:
	sh scripts/test-install.sh

vet:
	go vet ./...

diff-check:
	git diff --check
