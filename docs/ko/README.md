# investiq

이미 사용 중인 AI 도구로 투자 질문을 더 엄격하게 조사하세요.

`investiq`는 터미널에서 실행되는 작은 투자 리서치 앱입니다. 사용자의 질문을 현재 컴퓨터에서 사용할 수 있는 AI 도구에 자동으로 전달합니다. 사용할 수 있는 도구가 없으면 ChatGPT, Claude, Gemini 같은 웹 AI에 붙여넣을 수 있는 프롬프트를 만들어 줍니다.

English: [../../README.md](../../README.md)

## 설치

```bash
curl -fsSL https://raw.githubusercontent.com/dkstm95/investiq/main/install.sh | sh
```

설치하면 두 명령어가 생깁니다.

- `investiq`
- `iq`

앱을 열 때는 `investiq`를 쓰면 됩니다. 짧게 쓰고 싶으면 `iq`를 쓰면 됩니다.

Node, npm, Python, Go를 따로 설치할 필요는 없습니다.

설치 중 `export PATH=...` 문구가 나오면 터미널에 한 번 실행하세요.

## 시작하기

```bash
investiq
```

다음 화면이 열립니다.

```text
INVESTIQ
Fresh-data investment research for your AI tools

+------------------------------------------------------------+
| Ask a question. investiq will classify it automatically.   |
| It always asks the AI to use current web or market data.   |
+------------------------------------------------------------+

Ask >
```

질문을 입력하면 됩니다.

```text
Ask > NVDA 지금 사도 될까?
```

사용자가 모드나 AI 플랫폼을 고를 필요는 없습니다. `investiq`가 내부적으로 처리합니다.

## 질문 예시

```text
Ask > 6개월에서 12개월 관점에서 무엇에 투자하면 좋을까?
Ask > AI 인프라에 투자하고 싶은데 어떤 종목이나 ETF를 비교해야 할까?
Ask > TSLA를 보유 중인데 언제 줄이거나 팔아야 할까?
Ask > 내 포트폴리오를 점검해줘: AAPL 35%, NVDA 30%, SGOV 35%
```

앱을 열지 않고 한 번만 질문할 수도 있습니다.

```bash
iq "NVDA 지금 사도 될까?"
```

## 설정 확인

사용 가능한 AI 도구를 확인하려면 다음을 실행하세요.

```bash
iq doctor
```

`investiq`는 다음 순서로 사용할 도구를 찾습니다.

```text
Claude Code -> Codex -> OpenCode -> Gemini CLI -> 복사/붙여넣기 프롬프트
```

도구 실행에 실패하면 웹 AI에 붙여넣을 수 있는 프롬프트로 자동 전환합니다.

## 업데이트

```bash
iq update
```

현재 컴퓨터에 맞는 바이너리를 내려받아 기존 파일을 교체합니다.

`v0.1.4` 이하를 설치한 사용자는 업데이트 명령이 없으므로, 최초 1회는 설치 명령을 다시 실행해야 합니다.

```bash
curl -fsSL https://raw.githubusercontent.com/dkstm95/investiq/main/install.sh | sh
```

그 이후부터는 `iq update`를 사용할 수 있습니다.

## 최신 데이터

투자 정보는 빠르게 바뀝니다. 사용자가 "최신"이라고 쓰지 않아도 `investiq`는 AI에게 현재 웹/시장 데이터를 확인하라고 지시합니다.

현재 데이터를 확인할 수 없으면 AI는 추천을 내리지 않아야 합니다.

좋은 답변에는 다음이 포함되어야 합니다.

- 데이터 기준 시점
- 출처 링크
- 현재 가격 또는 최신 수치
- 핵심 아이디어
- 리스크
- 반대 의견
- 매수, 보유, 매도, 관망 조건
- 앞으로 볼 지표

## 한계

`investiq`는 다음을 하지 않습니다.

- 자동 매매
- 수익 보장
- 전문 금융, 세무, 법률 자문 대체
- 현재 버전에서 직접 시장 데이터 조회

`investiq`는 엄격한 리서치 흐름을 만들어 AI 도구에 전달합니다. 최신 웹/시장 데이터 조사는 연결된 AI 도구가 수행해야 합니다.

## 개발자 문서

개발 문서: [../development.md](../development.md)
