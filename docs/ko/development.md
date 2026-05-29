# 개발

## 소스에서 빌드

```bash
git clone https://github.com/dkstm95/gacha.git
cd gacha
go test ./...
go build -o gacha ./cmd/gacha
./gacha doctor
ln -sf gacha gch
./gch doctor
```

## 릴리즈

```bash
VERSION=0.1.12 sh scripts/build-release.sh
gh release create v0.1.12 dist/*.tar.gz dist/checksums.txt --title "v0.1.12"
```

GitHub Actions 템플릿은 다음 위치에 있습니다.

```text
docs/github-actions/
```

현재 repository token에 workflow 파일 push 권한이 없어서 템플릿으로 보관하고 있습니다.

## Codex 플러그인 파일

```text
.agents/plugins/marketplace.json
plugins/gacha/.codex-plugin/plugin.json
plugins/gacha/skills/gacha/SKILL.md
```
