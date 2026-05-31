.PHONY: check fmt test vet diff-check

export GOCACHE := $(CURDIR)/.cache/go-build
export GOMODCACHE := $(CURDIR)/.cache/go-mod

check: fmt test vet diff-check

fmt:
	mkdir -p $(GOCACHE) $(GOMODCACHE)
	gofmt -w ./cmd ./internal

test:
	go test ./...

vet:
	go vet ./...

diff-check:
	git diff --check
