# gacha

로컬 AI 런타임으로 투자 질문을 더 엄격하게 조사하세요.

`gacha`는 터미널에서 실행되는 투자 리서치 앱입니다. 내부 AI 런타임으로 OpenCode를 사용하므로, 사용자는 매번 플랫폼을 고르지 않고 ChatGPT, GitHub Copilot, Gemini, OpenAI API 또는 다른 지원 provider를 연결해 사용할 수 있습니다.

런타임이 없으면 첫 실행 시 `gacha`가 설치를 도와줍니다. 런타임을 사용할 수 없는 경우에도 ChatGPT, Claude, Gemini 같은 웹 AI에 붙여넣을 수 있는 프롬프트를 만들어 줍니다.

English: [../../README.md](../../README.md)

## 설치

```bash
curl -fsSL https://raw.githubusercontent.com/dkstm95/gacha/main/install.sh | sh
```

설치하면 기본 명령어와 짧은 별칭이 생깁니다.

- `gacha`
- `gch`

평소에는 짧은 `gch`를 쓰면 됩니다. 전체 이름이 필요할 때는 `gacha`를 쓰면 됩니다.

`gacha` 자체를 사용하기 위해 Node, npm, Python, Go를 따로 설치할 필요는 없습니다.

첫 실행 시 OpenCode runtime 설치와 AI provider 연결을 물어볼 수 있습니다. 이 runtime은 Gacha UI 뒤에서 AI를 실행하는 역할을 합니다.

설치 중 `export PATH=...` 문구가 나오면 터미널에 한 번 실행하세요.

## 시작하기

```bash
gch
```

다음 화면이 열립니다.

```text
GACHA
Fresh-data investment research for your AI tools

+------------------------------------------------------------+
| Ask a question. gacha will classify it automatically.      |
| It always asks the AI to use current web or market data.   |
+------------------------------------------------------------+

Ask >
```

질문을 입력하면 됩니다.

```text
Ask > NVDA 지금 사도 될까?
```

사용자가 모드나 AI 플랫폼을 고를 필요는 없습니다. `gacha`가 내부적으로 처리합니다.

## 질문 예시

```text
Ask > 6개월에서 12개월 관점에서 무엇에 투자하면 좋을까?
Ask > AI 인프라에 투자하고 싶은데 어떤 종목이나 ETF를 비교해야 할까?
Ask > TSLA를 보유 중인데 언제 줄이거나 팔아야 할까?
Ask > 내 포트폴리오를 점검해줘: AAPL 35%, NVDA 30%, SGOV 35%
```

앱을 열지 않고 한 번만 질문할 수도 있습니다.

```bash
gch "NVDA 지금 사도 될까?"
```

## 설정 확인

로컬 AI runtime 준비 상태를 확인하려면 다음을 실행하세요.

```bash
gch doctor
```

`gacha`는 다음 경로를 사용합니다.

```text
OpenCode runtime -> 복사/붙여넣기 프롬프트
```

OpenCode가 없거나 연결된 provider가 없다면 다음 명령을 실행하세요.

```bash
gch setup
```

`gch setup`은 필요한 경우 runtime을 설치한 뒤 provider login을 시작합니다. ChatGPT, GitHub Copilot, Gemini, OpenAI API 또는 OpenCode가 지원하는 다른 provider를 연결할 수 있습니다.

설정 후 투자 리서치 흐름과 결과 표시는 계속 Gacha UI 안에서 처리됩니다.

runtime 실행에 실패하면 웹 AI에 붙여넣을 수 있는 프롬프트로 자동 전환합니다.

## 업데이트

```bash
gch update
```

현재 컴퓨터에 맞는 바이너리를 내려받아 기존 파일을 교체합니다.

## 최신 데이터

투자 정보는 빠르게 바뀝니다. 사용자가 "최신"이라고 쓰지 않아도 `gacha`는 AI에게 현재 웹/시장 데이터를 확인하라고 지시합니다.

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

`gacha`는 다음을 하지 않습니다.

- 자동 매매
- 수익 보장
- 전문 금융, 세무, 법률 자문 대체
- 현재 버전에서 직접 시장 데이터 조회

`gacha`는 엄격한 리서치 흐름을 만들어 AI 도구에 전달합니다. 최신 웹/시장 데이터 조사는 연결된 AI 도구가 수행해야 합니다.

## 개발자 문서

개발 문서: [../development.md](../development.md)
