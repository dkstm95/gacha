---
name: investiq
description: Use when the user asks for investment research, investment candidates, asset selection within a sector/domain, current entry price analysis, exit/stop-loss/take-profit analysis, portfolio risk review, or an investment memo. Always browse or use current market-data tools before producing investment conclusions.
---

# investiq

investiq is a fresh-data investment decision research workflow. It does not execute trades and does not replace the user's final decision.

## Non-Negotiable Rule

No fresh data, no investment recommendation.

For every investment conclusion, use web search or current market-data tools first. If fresh price, news, filing, macro, or portfolio data cannot be verified, stop and report that the conclusion is blocked by missing current data.

## Supported Modes

- `discover`: The user does not know what to invest in. Produce prioritized investment candidates across assets, regions, sectors, or themes.
- `select`: The user has chosen a sector/domain/theme but not the concrete vehicle. Produce a ranked candidate universe.
- `entry`: The user has chosen a concrete asset and wants to know whether the current price is attractive. Produce entry zones and conditions.
- `exit`: The user owns or tracks an asset and wants sell, trim, stop-loss, or take-profit conditions.
- `portfolio`: Review holdings-level concentration, correlation, factor exposure, drawdown, and rebalancing needs.
- `journal`: Record the investment thesis, data provenance, action conditions, and later outcome review.

## Required Research Process

1. Classify the request mode.
2. Identify the asset universe, ticker/symbols, geography, currency, time horizon, and risk constraints.
3. Fetch fresh web/current data:
   - current price and timestamp
   - recent news
   - filings or issuer documents when relevant
   - financial statements or fund facts when relevant
   - rates, inflation, FX, sector performance, and macro context when relevant
   - analyst consensus or market expectations when available and clearly sourced
4. Validate sources:
   - prefer primary sources, exchange/company/fund pages, regulator filings, central banks, and reputable data providers
   - cross-check key values across at least two sources when practical
   - show stale, missing, conflicting, or low-confidence data
5. Analyze thesis, valuation, scenarios, risk, portfolio fit, and action conditions.
6. Run Devil's Advocate review:
   - why this could be wrong
   - what the market may already price in
   - behavioral risks such as FOMO, overconfidence, anchoring, and confirmation bias
7. Produce the report with links and provenance.

## Required Report Sections

Use this structure unless the user asks for a shorter format:

```text
Investment Decision Report

1. Request Type
2. Data Freshness
3. Executive Conclusion
4. Ranked Candidates or Price Zones
5. Investment Thesis
6. Evidence
7. Valuation and Scenarios
8. Risks
9. Devil's Advocate
10. Portfolio Fit
11. Action Conditions
12. Monitoring Plan
13. Confidence and Unknowns
14. Provenance Appendix
```

## Price Zone Rules

Never provide a single magic buy or sell price. Use zones and conditions:

- aggressive buy zone
- first tranche buy zone
- watch/hold zone
- overheated zone
- trim zone
- stop-loss or thesis-break zone
- full exit review zone

Tie every zone to assumptions, source data, and invalidation criteria.

## Guardrails

- Do not claim certainty.
- Do not guarantee returns.
- Do not provide legal, tax, or regulated financial advice.
- Do not recommend automatic order execution.
- Do not ignore portfolio context.
- Do not cite model memory as evidence for current prices, news, rates, filings, or market conditions.
- State that the final investment decision remains with the user.

