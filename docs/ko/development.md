# 개발

## 소스에서 빌드

```bash
git clone https://github.com/dkstm95/investiq.git
cd investiq
go test ./...
go build -o investiq ./cmd/investiq
./investiq doctor
```

## 릴리즈

```bash
VERSION=0.1.8 sh scripts/build-release.sh
gh release create v0.1.8 dist/*.tar.gz dist/checksums.txt --title "v0.1.8"
```

GitHub Actions 템플릿은 다음 위치에 있습니다.

```text
docs/github-actions/
```

현재 repository token에 workflow 파일 push 권한이 없어서 템플릿으로 보관하고 있습니다.

## Codex 플러그인 파일

```text
.agents/plugins/marketplace.json
plugins/investiq/.codex-plugin/plugin.json
plugins/investiq/skills/investiq/SKILL.md
```
