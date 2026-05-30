.PHONY: check fmt test vet diff-check

check: fmt test vet diff-check

fmt:
	gofmt -w ./cmd ./internal

test:
	go test ./...

vet:
	go vet ./...

diff-check:
	git diff --check
