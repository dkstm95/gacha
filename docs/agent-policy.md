# Agent Policy

This document is the source of truth for AI agents working on this repository.
Agent-specific files, such as `AGENTS.md`, should reference this file instead
of duplicating policy.

## Quality Pass

Before finishing any code change, review the changed diff for:

- bugs or behavioral regressions
- duplicated logic or unclear ownership
- overengineering or unnecessary features
- confusing user flow, copy, or terminal UI polish issues
- missing tests for changed behavior

Keep follow-up refactors small and directly related to the touched area. Prefer
deleting, moving, or simplifying code over adding new abstractions. Do not add
features unless the user explicitly asks.

## Product UI Philosophy

For UI work, follow the product-level guidance in:

```text
docs/principles.md
```

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

## Git Automation

After requested work is complete and verification passes, the agent may commit
without asking the user for a separate decision.

Before committing:

- inspect `git status --short`
- include only files related to the completed task
- do not include unrelated user changes
- do not commit if verification failed or could not be run
- use a concise commit message that describes the completed change

If unrelated changes are present, leave them unstaged and mention them in the
final response.

## Release Automation

After a successful commit, the agent may decide whether to release without
asking the user for a separate decision.

Release only when all of these are true:

- the change is user-facing, affects distributed behavior, or updates release
  artifacts or installation behavior
- `make check` passed after the final diff
- the working tree contains no unrelated staged changes
- release credentials and required tools are already available
- the repository's release process can be completed without manual secrets,
  account setup, or interactive approval

Use the smallest appropriate semantic version bump:

- patch for fixes, UI polish, docs corrections, and low-risk behavior changes
- minor for new user-facing features
- major only for intentional breaking changes

Do not release when the change is only local agent policy, test-only cleanup,
or repository maintenance that does not affect users.

Never publish a release if build, test, packaging, tag, push, or release
creation fails. Report the failure and leave the repository in the clearest
recoverable state.
