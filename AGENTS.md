# Gacha Agent Instructions

Use these repository-specific rules after every code change.

## Quality Pass

Before the final response, review the changed diff for:

- bugs or behavioral regressions
- duplicated logic or unclear ownership
- overengineering or unnecessary features
- confusing user flow, copy, or terminal UI polish issues
- missing tests for changed behavior

Keep follow-up refactors small and directly related to the touched area. Prefer deleting, moving, or simplifying code over adding new abstractions. Do not add features unless the user explicitly asks.

## Verification

Run the standard project check after code changes:

```sh
make check
```

If `make` is unavailable, run the equivalent commands:

```sh
gofmt -w ./cmd ./internal
go test ./...
go vet ./...
git diff --check
```

Mention any verification that could not be run.

## Release Safety

Do not commit, push, bump versions, tag releases, or publish releases unless the user explicitly asks.
