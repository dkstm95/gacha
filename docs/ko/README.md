# investiq

`investiq`는 standalone 투자 리서치 AI agent harness입니다.

사용자의 투자 질문을 현재 사용할 수 있는 AI CLI에 자동으로 전달하고, 항상 최신 데이터 조사, 출처 링크, 리스크 검토, 반대 논리, provenance를 포함하도록 강제합니다. 사용할 수 있는 AI CLI가 없으면 웹 검색 가능한 AI에 붙여넣을 수 있는 프롬프트를 출력합니다.

## 설치

```bash
curl -fsSL https://raw.githubusercontent.com/dkstm95/investiq/main/install.sh | sh
```

설치하면 두 명령어가 생깁니다.

- `investiq`
- `iq`

일반 사용자는 짧은 `iq`를 쓰면 됩니다.

설치 중 PATH 안내가 나오면 출력된 `export PATH=...` 명령을 실행하세요.

## 빠른 시작

```bash
iq init
iq doctor
investiq
```

`investiq`를 실행하면 interactive UI로 진입합니다.

```text
investiq
Fresh-data investment research agent

Ask an investment question. investiq will classify it and route it automatically.
Type /help for commands, /doctor to check AI platforms, /quit to exit.

iq>
```

`iq>` 프롬프트에 질문을 입력하면 됩니다.

```text
iq> NVDA 지금 사도 될까?
```

사용자가 `entry`, `exit`, `platform` 같은 옵션을 고를 필요가 없습니다. `investiq`가 내부적으로 요청 종류를 분류하고 사용할 AI 플랫폼을 자동 선택합니다.

예시:

```text
iq> 6개월에서 12개월 관점에서 무엇에 투자하면 좋을까?
iq> AI 인프라에 투자하고 싶은데 어떤 종목이나 ETF를 비교해야 할까?
iq> TSLA를 보유 중인데 언제 줄이거나 팔아야 할까?
iq> 내 포트폴리오를 점검해줘: AAPL 35%, NVDA 30%, SGOV 35%
```

한 번만 실행하고 싶으면 기존처럼 짧게 쓸 수도 있습니다.

```bash
iq "NVDA 지금 사도 될까?"
```

## 동작 방식

`investiq`는 로컬에서 사용 가능한 AI CLI를 확인하고 다음 순서로 자동 선택합니다.

```text
Claude Code -> Codex -> OpenCode -> Gemini CLI -> manual prompt
```

감지된 플랫폼이 실행 중 실패하면 사용자가 막히지 않도록 자동으로 프롬프트 출력 방식으로 fallback합니다.

`iq init`은 다음 파일을 생성합니다.

```text
~/.investiq/config.json
```

라우팅 순서나 명령어 이름을 바꾸고 싶으면 이 파일을 수정하면 됩니다.

## 최신 데이터 원칙

사용자가 "최신", "현재", "최근" 같은 표현을 쓰지 않아도 모든 투자 분석은 최신 웹/시장 데이터를 조사해야 합니다.

최신 데이터를 검증할 수 없으면 추천을 내리지 않아야 합니다.

결과에는 다음이 포함되어야 합니다.

- 데이터 기준 시점
- 출처 링크
- 현재 가격 또는 최신 수치
- 투자 thesis
- 밸류에이션 또는 시나리오 분석
- 리스크
- 반대 논리
- 행동 조건
- 모니터링 계획
- provenance appendix

## 하지 않는 것

`investiq`는 다음을 하지 않습니다.

- 자동 매매 실행
- 수익 보장
- 전문 금융/세무/법률 자문 대체
- 현재 버전에서 직접 시장 데이터 API 호출

현재 버전은 엄격한 투자 리서치 워크플로우를 구성하고 AI 플랫폼에 전달하는 역할을 합니다. 최신 데이터 조사는 연결된 AI 플랫폼이 수행해야 합니다.

## 설계 메모

초기 설계와 리서치 노트는 [design-notes.md](design-notes.md)에 있습니다.
