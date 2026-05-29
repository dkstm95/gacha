# gacha

`gacha` is a marketplace-ready investment research agent harness. The default product direction is a Gacha-owned terminal UI backed by a local OpenCode runtime. Users connect ChatGPT, GitHub Copilot, Gemini, API providers, or other OpenCode-supported providers through that runtime instead of choosing a platform on every request.

## Quick Start

```bash
curl -fsSL https://raw.githubusercontent.com/dkstm95/gacha/main/install.sh | sh
gacha setup
gch doctor
gch "NVDA 지금 사도 될까?"
```

The installer downloads a standalone binary from GitHub Releases and installs the `gacha` command and the `gch` short alias. It does not require Node, npm, Python, or Go for the Gacha binary. On first run, Gacha can install OpenCode runtime and start provider login for the user.

Install a specific release:

```bash
curl -fsSL https://raw.githubusercontent.com/dkstm95/gacha/main/install.sh | GACHA_VERSION=v0.1.9 sh
```

Build from source:

```bash
git clone https://github.com/dkstm95/gacha.git
cd gacha
go build -o gacha ./cmd/gacha
./gacha doctor
```

Maintainer release flow:

```bash
VERSION=0.1.9 sh scripts/build-release.sh
gh release create v0.1.9 dist/*.tar.gz dist/checksums.txt --title "v0.1.9"
```

Codex marketplace plugin:

```text
.agents/plugins/marketplace.json
plugins/gacha/.codex-plugin/plugin.json
```

Embedded agent assets:

```text
internal/app/assets/plugins/gacha/platforms/generic/system-prompt.md
internal/app/assets/plugins/gacha/workflows/
internal/app/assets/plugins/gacha/templates/
```

The product name and full command are `gacha`. The day-to-day short command is `gch`.

The CLI composes host-agnostic workflows and templates, then routes them to OpenCode runtime. It does not fetch market data or execute trades by itself. The connected AI provider must use current web or market-data tools before producing investment conclusions.

Harness commands:

```bash
gch doctor
gch setup
gch "AAPL 지금 살까?"
gch "AAPL 현재 매수 구간 분석"
gch "TSLA 보유 중인데 매도 기준 점검"
```

Users connect their actual subscriptions through `gch setup`, which delegates credential storage to OpenCode. Gacha keeps the investment workflow and UI on top. The runtime route is intentionally fixed:

```text
OpenCode runtime → copy/paste prompt fallback
```

## 1. Project Purpose

`gacha`는 투자 판단을 돕는 AI agent 팀을 설계하기 위한 프로젝트이다. 프로젝트명은 이후 CLI 명령어 이름으로도 사용한다.

목표는 AI가 사용자를 대신해 매수/매도 결정을 내리는 것이 아니라, 최신 데이터와 신뢰 가능한 출처를 기반으로 투자 후보를 좁히고, 매수/매도 가격대와 리스크를 체계적으로 검토하도록 돕는 것이다.

핵심 원칙은 다음과 같다.

```text
AI Investment Decision Team =
최신 데이터 기반 투자 리서치
+ 후보 우선순위화
+ 매수/매도 가격대 분석
+ 리스크와 반대 논리 검증
+ 출처 기반 투자 메모 생성
```

## 2. Core Operating Rule

투자 관련 정보는 가격, 실적, 금리, 뉴스, 공시, 규제, ETF 구성, 밸류에이션이 모두 빠르게 바뀐다. 따라서 이 팀은 항상 웹 서치를 통해 사용자의 요청 시점 이후 확인 가능한 최신 데이터를 기반으로 조사하고 분석해야 한다.

```text
No fresh data, no investment recommendation.
```

웹 서치를 못 했거나 최신 가격/뉴스/실적/공시를 확인하지 못했다면, 추천을 내리지 않고 "최신 데이터 확인 불가로 결론 보류"라고 보고해야 한다.

모든 결과물에는 다음 항목이 포함되어야 한다.

- 사용한 데이터의 기준일
- 주요 출처 링크
- 현재 가격 또는 최신 수치
- 데이터 제공자와 조회 시각
- 데이터가 불확실한 부분
- 추가 확인이 필요한 지표
- 결론의 신뢰도

또한 모든 핵심 수치에는 provenance를 남긴다.

```text
Provenance =
source name
+ source URL
+ retrieved_at
+ data_as_of
+ symbol / identifier
+ raw value
+ normalized value
+ confidence
```

## 3. Target Use Cases

### 3.1 투자 대상을 모를 때

사용자가 무엇에 투자해야 할지 모르는 경우, AI 팀은 거시경제, 시장 흐름, 산업/섹터, 자산군을 조사한 뒤 투자 후보들을 우선순위로 제안한다.

분석 대상 예시는 다음과 같다.

- 주식
- ETF
- 채권
- 원자재
- 부동산/REITs
- 크립토
- 현금성 자산
- 특정 국가/지역/섹터

### 3.2 섹터나 도메인은 정했지만 구체적 대상을 모를 때

예를 들어 사용자가 "AI 인프라에 투자하고 싶다", "반도체 섹터에 투자하고 싶다", "미국 배당주를 찾고 싶다"처럼 말하면, AI 팀은 해당 영역의 투자 universe를 만들고 후보를 비교한다.

비교 기준은 다음과 같다.

- 성장성
- 수익성
- 밸류에이션
- 경쟁력
- 재무 안정성
- 모멘텀
- 리스크
- 포트폴리오 적합성

### 3.3 구체적 투자 대상의 매수 시점이 궁금할 때

사용자가 특정 종목, ETF, 코인, 채권, 원자재 등을 정했지만 현재 가격이 매수하기 적절한지 모를 때, AI 팀은 최신 가격과 데이터를 조사해 적정 매수 가격대를 제안한다.

출력은 단일 가격보다 구간 중심이어야 한다.

```text
공격적 매수 구간
1차 분할매수 구간
관망 구간
과열 구간
```

### 3.4 보유 자산의 매도 시점이 궁금할 때

사용자가 이미 투자 중인 자산을 언제 매도해야 할지 모를 때, AI 팀은 최신 데이터를 기반으로 손절, 부분 익절, 전량 매도 검토 구간을 제안한다.

매도 판단은 세 가지 기준으로 나눈다.

```text
1. Thesis-based Exit
   투자 가설이 틀렸을 때 매도

2. Price-based Exit
   목표가, 손절가, 과열 구간 도달 시 매도

3. Portfolio-based Exit
   비중 과대, 더 좋은 기회, 리밸런싱 필요 시 매도
```

## 4. Recommended Agent Team

### 4.1 Investment Policy Agent

사용자의 투자 원칙과 제약을 관리한다.

책임:

- 투자 목적 확인
- 투자 기간 확인
- 위험 허용도 확인
- 선호/비선호 자산군 확인
- 유동성, 세금, 통화, 국가 제한 확인
- 후보가 사용자의 투자 원칙에 맞는지 1차 필터링

### 4.2 Web Research Agent

최신 웹 데이터를 수집한다.

책임:

- 현재 가격 확인
- 최근 뉴스 확인
- 실적 발표, 공시, 리포트 확인
- 금리, 환율, 인플레이션, 경기 지표 확인
- 섹터/산업 트렌드 확인
- ETF 구성 및 비용 확인

### 4.3 Source Validator Agent

출처와 데이터 품질을 검증한다.

책임:

- 출처의 신뢰도 확인
- 데이터 기준일 확인
- 서로 다른 출처 간 수치 비교
- 오래된 정보와 최신 정보 구분
- 불확실한 데이터 표시
- 핵심 수치의 provenance 기록
- 가격, 재무, 뉴스, 공시 데이터의 source mismatch 표시

### 4.4 Opportunity Scout Agent

투자 기회를 발굴한다.

책임:

- 자산군별 투자 환경 조사
- 섹터/테마별 기회 탐색
- 후보군 생성
- 투자 매력도와 리스크의 1차 비교

### 4.5 Asset Selection Agent

구체적인 투자 대상을 비교하고 우선순위를 정한다.

책임:

- 동일 섹터/도메인 내 후보 비교
- 주식, ETF, 채권, 원자재, 크립토 등 후보별 장단점 정리
- 퀄리티, 밸류에이션, 성장성, 모멘텀, 리스크 점수화
- 최종 후보 랭킹 작성

### 4.6 Valuation & Entry Agent

매수 가격대와 진입 전략을 분석한다.

책임:

- 현재 가격과 과거 밸류에이션 비교
- forward valuation 검토
- base/bull/bear scenario 작성
- upside/downside 계산
- 분할매수 구간 제안
- 매수 보류 조건 제안

### 4.7 Exit & Risk Agent

매도, 손절, 리밸런싱 기준을 분석한다.

책임:

- 손절 가격대 제안
- 부분 익절 구간 제안
- 전량 매도 검토 조건 제안
- thesis invalidation 조건 정의
- 포트폴리오 집중도와 상관관계 점검
- 최대 손실 가능성 검토

### 4.8 Devil's Advocate Agent

반대 논리와 행동편향을 점검한다.

책임:

- 투자 thesis의 약점 공격
- 과신, 확증편향, FOMO, 군중심리 점검
- bear case 강화
- "이 투자가 틀릴 수 있는 이유" 제시
- 추천 보류 또는 반대 의견 작성

### 4.9 Investment Committee Agent

최종 투자 메모를 작성한다.

책임:

- 각 agent의 분석 취합
- 우선순위 또는 가격대 결정
- 근거, 리스크, 반대 논리 정리
- 최종 결론을 Buy / Watch / Avoid / Trim / Sell 후보 등으로 구조화
- 모니터링 계획과 재검토 시점 제안

### 4.10 Trade Journal Agent

투자 판단과 사후 결과를 기록한다.

책임:

- 최초 투자 thesis 기록
- 매수/매도 판단 당시의 데이터와 출처 저장
- 사후 성과와 판단 품질 비교
- 틀린 가정과 개선 포인트 기록
- 반복되는 행동편향 탐지

MVP에서는 독립 agent가 아니라 Investment Committee Agent의 하위 기능으로 시작해도 된다.

## 5. Workflows

### 5.1 Discover Mode

사용자가 무엇에 투자할지 모를 때 사용한다.

```text
User Request
→ Investment Policy Agent
→ Web Research Agent
→ Source Validator Agent
→ Opportunity Scout Agent
→ Exit & Risk Agent
→ Devil's Advocate Agent
→ Investment Committee Agent
```

출력:

- 투자 후보 우선순위
- 후보별 투자 논리
- 기대 수익 요인
- 주요 리스크
- 적합한 투자 기간
- 추천 진입 방식
- 재검토 시점

### 5.2 Select Mode

사용자가 섹터/도메인은 정했지만 구체적 투자 대상을 모를 때 사용한다.

```text
User Request
→ Investment Policy Agent
→ Web Research Agent
→ Source Validator Agent
→ Opportunity Scout Agent
→ Asset Selection Agent
→ Valuation & Entry Agent
→ Devil's Advocate Agent
→ Investment Committee Agent
```

출력:

- 후보 universe
- 후보별 비교표
- 투자 매력도 점수
- 리스크 점수
- 최종 우선순위
- 추천/보류/제외 사유

### 5.3 Entry Mode

사용자가 구체적 투자 대상의 현재 매수 적절성을 알고 싶을 때 사용한다.

```text
User Request
→ Investment Policy Agent
→ Web Research Agent
→ Source Validator Agent
→ Valuation & Entry Agent
→ Exit & Risk Agent
→ Devil's Advocate Agent
→ Investment Committee Agent
```

출력:

- 현재 가격
- 공격적 매수 구간
- 1차 분할매수 구간
- 관망 구간
- 과열 구간
- 핵심 가정
- 주요 리스크
- thesis invalidation 조건

### 5.4 Exit Mode

사용자가 보유 자산의 매도 시점, 손절, 익절 기준을 알고 싶을 때 사용한다.

```text
User Request
→ Investment Policy Agent
→ Web Research Agent
→ Source Validator Agent
→ Exit & Risk Agent
→ Valuation & Entry Agent
→ Devil's Advocate Agent
→ Investment Committee Agent
```

출력:

- 손절/축소 구간
- 부분 익절 구간
- 전량 매도 검토 조건
- thesis-based exit 조건
- portfolio-based exit 조건
- 계속 보유할 조건
- 다음 모니터링 이벤트

## 6. Standard Report Format

모든 결과물은 다음 형식을 따른다.

```text
Investment Decision Report

1. Request Type
   Discover / Select / Entry / Exit

2. Data Freshness
   사용한 데이터 기준일, 가격 기준 시각, 주요 출처

3. Executive Conclusion
   핵심 결론과 추천 등급

4. Ranked Candidates or Price Zones
   후보 우선순위 또는 매수/매도 가격대

5. Investment Thesis
   왜 이 투자 아이디어가 성립하는가

6. Evidence
   최신 데이터, 공시, 실적, 뉴스, 거시 지표

7. Valuation and Scenarios
   base / bull / bear case

8. Risks
   핵심 리스크와 손실 가능성

9. Devil's Advocate
   반대 논리와 실패 시나리오

10. Portfolio Fit
   기존 포트폴리오와의 관계, 집중도, 상관관계

11. Action Conditions
   매수, 관망, 손절, 익절, 전량 매도 조건

12. Monitoring Plan
   추적할 지표, 이벤트, 재검토 시점

13. Confidence and Unknowns
   결론 신뢰도와 확인하지 못한 정보

14. Provenance Appendix
   주요 수치별 출처, 조회 시각, 데이터 기준일
```

## 7. Evaluation and Quality Gates

투자 AI 팀은 그럴듯한 설명을 생성하는 것보다 판단 품질을 반복적으로 개선하는 것이 중요하다. 따라서 `gacha`는 다음 평가 기준을 가져야 한다.

### 7.1 Report Quality

- 최신 데이터 사용 여부
- 출처 다양성
- 핵심 수치의 cross-check 여부
- thesis, risk, invalidation criteria 포함 여부
- 반대 논리의 구체성
- 가격대와 조건의 명확성

### 7.2 Decision Quality

- 추천 당시의 기대 수익과 실제 결과 비교
- bull/base/bear scenario의 적중 범위
- 손절/익절/관망 조건의 유효성
- 포트폴리오 집중도 변화
- drawdown 억제 여부

### 7.3 Process Quality

- 데이터 기준일 누락 여부
- 출처 링크 누락 여부
- 오래된 데이터 사용 여부
- 단일 출처 의존 여부
- 사용자의 투자 정책과 충돌 여부
- 자동매매처럼 보이는 표현 사용 여부

### 7.4 Backtesting and Paper Trading

후속 버전에서는 후보 랭킹, entry zone, exit rule을 과거 데이터로 검증할 수 있어야 한다. 단, 백테스트는 과최적화와 생존자 편향에 취약하므로 실제 추천의 유일한 근거로 사용하지 않는다.

필수 주의점:

- survivorship bias 점검
- look-ahead bias 방지
- 거래비용과 세금 반영
- 분할매수/분할매도 rule 반영
- 시장 국면별 성과 분리

### 7.5 Decision Journal

모든 투자 판단은 나중에 검토할 수 있도록 기록한다.

기록 항목:

- 요청 원문
- 사용한 데이터 기준일
- 핵심 출처
- 추천 당시 가격
- 투자 thesis
- 반대 논리
- action condition
- 사후 결과
- 판단 오류와 개선점

## 8. Design Principles from Research

### 8.0 유사 프로젝트와 제품에서 확인한 패턴

웹 서치로 확인한 유사 오픈소스, 제품, 논문/연구에서 반복적으로 나타난 패턴은 다음과 같다.

1. 멀티에이전트 구조는 유효하지만 역할이 너무 많으면 운영 비용이 커진다.
2. 최신 시장 데이터, 재무 데이터, 뉴스, 공시를 가져오는 도구 계층이 핵심이다.
3. 금융 분석에서는 "모델의 설명"보다 "검증 가능한 데이터 provenance"가 더 중요하다.
4. 좋은 시스템은 추천만 하지 않고 리스크, 반대 논리, 포트폴리오 맥락, 모니터링 조건을 함께 제공한다.
5. 자동매매보다 리서치 보조, 보고서 생성, 포트폴리오 점검, 알림, 투자 저널이 MVP에 더 적합하다.

대표 사례:

- TradingAgents: fundamental analyst, sentiment analyst, technical analyst, bull/bear researcher, trader, risk manager 등으로 trading firm을 모사하는 multi-agent framework.
- FinRobot: 금융 문제를 Financial Chain-of-Thought로 분해하고, agent layer, model strategy layer, LLMOps/DataOps layer, multi-source model layer로 구성한 오픈소스 금융 agent 플랫폼.
- FinGPT: 금융 특화 LLM, sentiment, forecasting, retrieval, multi-agent analysis, report intelligence, benchmarking을 포함하는 오픈소스 금융 AI 생태계.
- TradeApe: agent-first chart analysis를 지향하며, time-series database와 inspectable tools, visible data provenance를 강조하는 오픈소스 프로젝트.
- Financial Research Analyst Agent: 11개 specialized agents, 20개 이상 분석 도구, RAG pipeline, multi-provider data layer, CLI/API/UI를 제공하는 오픈소스 stock analysis 시스템.
- FullAriza: 10개 specialized AI agents, algorithmic signals, portfolio x-ray, VaR/CVaR, backtesting, tax optimization, alerts를 제공하는 상용 제품.
- MSCI AI Portfolio Insights: natural-language portfolio risk analysis를 지원하며, agent가 context 추가, 포트폴리오 데이터 조립, 코드 생성으로 risk query에 답하는 방식을 설명한다.

References:

- https://arxiv.org/abs/2412.20138
- https://tradingagents.co/
- https://arxiv.org/abs/2405.14767
- https://fingpt.io/
- https://tradeape.org/
- https://github.com/gsaini/financial-research-analyst-agent
- https://fullariza.com/en/product
- https://www.msci.com/research-and-insights/paper/ai-portfolio-insights-and-the-future-of-risk-management

### 8.1 투자 프로세스는 정책에서 시작한다

CFA Institute의 포트폴리오 관리 프레임워크는 투자자의 목표와 제약을 먼저 이해하고, Investment Policy Statement를 만든 뒤 자산배분, 종목 분석, 포트폴리오 구성, 모니터링, 리밸런싱, 성과 측정으로 이어지는 과정을 강조한다.

AI 팀도 개별 아이디어보다 먼저 사용자의 투자 목적, 기간, 위험 허용도, 제약 조건을 확인해야 한다.

Reference:

- https://www.cfainstitute.org/insights/professional-learning/refresher-readings/2026/portfolio-management-overview

### 8.2 개별 투자 아이디어보다 포트폴리오 전체가 중요하다

BlackRock과 Vanguard의 포트폴리오 구성 원칙은 명확한 목표, 자산배분, 비용, 리밸런싱, 장기 규율을 강조한다.

AI 팀은 "이 종목이 좋은가"만 판단하지 않고, "이 투자가 전체 포트폴리오에서 어떤 역할을 하는가"를 함께 분석해야 한다.

References:

- https://www.blackrock.com/americas-offshore/en/education/portfolio-construction/understanding-portfolio-construction
- https://corporate.vanguard.com/content/dam/corp/research/pdf/vanguards_principles_for_investing_success.pdf

### 8.3 반대 의견은 구조적으로 포함되어야 한다

Bridgewater의 idea meritocracy는 독립적인 사고, 반대 의견, 근거 기반 토론을 통해 의사결정 품질을 높이는 방식을 강조한다.

AI 팀에는 반드시 Devil's Advocate Agent가 포함되어야 하며, 모든 추천은 반대 논리를 통과해야 한다.

References:

- https://www.bridgewater.com/culture/bridgewaters-idea-meritocracy
- https://www.principles.com/principles/633d5d13-8610-425f-ad62-cd62347d9165

### 8.4 행동편향을 통제해야 한다

투자는 확증편향, 과신, 손실회피, FOMO, 군중심리에 취약하다. Michael Mauboussin은 좋은 투자 프로세스를 analytical, behavioral, organizational 요소로 나누어 설명한다.

AI 팀은 분석뿐 아니라 행동편향 점검, 의사결정 기록, 사후 검증을 포함해야 한다.

Reference:

- https://www.fool.com/investing/general/2014/03/02/the-three-components-of-a-great-investment-process.aspx

### 8.5 금융 AI에는 거버넌스와 검증이 필요하다

Deloitte는 금융권 agentic AI가 agent owner, validator, steward와 같은 통제 구조를 가져야 하며, 감사 가능성과 인간 승인 절차가 필요하다고 설명한다.

AI 팀은 자동 매매보다 리서치, 모니터링, 리스크 검토, 투자 메모 생성부터 시작해야 한다.

Reference:

- https://www.deloitte.com/us/en/insights/industry/financial-services/agentic-ai-risks-banking.html

### 8.6 `gacha`의 차별화 방향

유사 프로젝트 중 상당수는 trading framework, stock report generator, portfolio analytics product에 가깝다. `gacha`는 다음 위치를 목표로 한다.

```text
gacha =
개인 투자자를 위한 최신 데이터 기반 investment decision copilot
```

따라서 초기 MVP는 자동매매나 고빈도 트레이딩보다 다음 기능에 집중한다.

- 투자 후보 발견
- 특정 도메인 내 후보 우선순위화
- 매수 가격대 분석
- 매도/손절/익절 조건 분석
- 투자 thesis와 반대 논리 기록
- 데이터 출처와 판단 근거 저장

하지 않는 것:

- 자동 주문 실행
- 수익 보장형 종목 추천
- 검증되지 않은 단일 데이터 출처 기반 결론
- 설명 불가능한 black-box score만 제공

## 9. Guardrails

이 팀은 다음 규칙을 지켜야 한다.

- 최신 웹 데이터 없이 투자 추천 금지
- 출처 링크 없는 핵심 주장 금지
- provenance 없는 핵심 수치 사용 금지
- 단일 가격만 제시하지 말고 가격대와 조건 제시
- 항상 반대 논리 포함
- 항상 손실 가능성과 thesis invalidation 조건 포함
- 포트폴리오 맥락 없이 개별 자산 추천 금지
- 자동 매매 실행 금지
- 최종 결정은 사용자가 내린다는 점 명시
- 백테스트 결과를 미래 수익 보장처럼 표현 금지

## 10. Future Extensions

추후 확장 가능한 agent는 다음과 같다.

- Macro Analyst Agent
- Quant Analyst Agent
- Technical Analysis Agent
- Options / Derivatives Analyst Agent
- Crypto Analyst Agent
- Real Estate Analyst Agent
- Tax / Regulation Watcher Agent
- News Monitoring Agent
- Portfolio Rebalancing Agent
- Trade Journal Agent
- Alerting Agent
- Backtesting Agent
- Data Connector Agent

## 11. MVP Implementation Direction

처음에는 복잡한 명령 모드보다 단일 질문 인터페이스로 시작한다. Gacha가 요청을 내부적으로 다음 유형 중 하나로 분류한다.

```text
discover
무엇에 투자할지 모를 때 후보 추천

select
섹터/도메인 안에서 구체적 투자 대상 추천

entry
특정 대상의 매수 적정 가격대 분석

exit
보유 자산의 매도/손절/익절 기준 분석
```

각 명령은 반드시 다음 순서를 따른다.

```text
1. 요청 분류
2. 투자 정책 확인
3. 최신 웹 데이터 수집
4. 출처 검증
5. 분석
6. 반대 논리 검토
7. 투자 메모 작성
8. 사용자 최종 판단 지원
```

MVP 이후에는 다음 내부 workflow를 추가할 수 있다.

```text
journal
투자 판단과 사후 결과 기록

monitor
보유 자산의 가격, 뉴스, 공시, thesis invalidation 조건 모니터링

portfolio
포트폴리오 집중도, 상관관계, drawdown, 리밸런싱 필요성 점검

backtest
후보 랭킹, entry zone, exit rule의 과거 성과 검증
```
