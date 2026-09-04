# OpenAPI Generation — Kanban Board

Epic: **OpenAPI schema generation via outbox + worker** — see issue #11.
Milestone: **OpenAPI Generation** (due 2026-09-30) — see milestone #1.

A user can request an OpenAPI 3.0 schema derived from the database records. The
request is enqueued through a transactional outbox, processed by a background
worker, and the finished document is stored and downloaded by polling
`GET /generations/:id`.

> Source of truth: the GitHub Issues tab. This file mirrors it. Move a ticket by
> changing its `status:` label (`status:todo` / `status:in-progress` /
> `status:done` / `status:blocked`).

## 🟦 Epic: OpenAPI schema generation

| Ticket | Title | Status | Priority |
| ------ | ----- | :----: | :------: |
| [#2](https://github.com/frisboo-bank/openapi-generator/issues/2) | Storage & queue foundation | ⬜ Todo | High |
| [#3](https://github.com/frisboo-bank/openapi-generator/issues/3) | OpenAPI 3.0 schema builder | ⬜ Todo | High |
| [#4](https://github.com/frisboo-bank/openapi-generator/issues/4) | Request API + background worker | ⬜ Todo | High |
| [#5](https://github.com/frisboo-bank/openapi-generator/issues/5) | End-to-end integration test | ⬜ Todo | Medium |

## 🐛 Bug fixes (carry alongside the feature)

| Ticket | Title | Status | Priority |
| ------ | ----- | :----: | :------: |
| [#6](https://github.com/frisboo-bank/openapi-generator/issues/6) | Working tree does not compile (build blockers) | ⬜ Todo | High |
| [#7](https://github.com/frisboo-bank/openapi-generator/issues/7) | Repair the metrics subsystem | ⬜ Todo | High |
| [#8](https://github.com/frisboo-bank/openapi-generator/issues/8) | Complete or remove the tracer subsystem | ⬜ Todo | Medium |
| [#9](https://github.com/frisboo-bank/openapi-generator/issues/9) | Entity repository read paths + Redis adapter | ⬜ Todo | High |
| [#10](https://github.com/frisboo-bank/openapi-generator/issues/10) | sql_client resolve hook blocks on Ping | ⬜ Todo | Medium |

## Legend

| Label | Meaning |
| ----- | ------- |
| `status:todo` | Not yet started |
| `status:in-progress` | Currently being worked on |
| `status:done` | Finished |
| `status:blocked` | Waiting on something external |
| `priority:high` | Do this next |
| `priority:medium` | Scheduled soon |
| `priority:low` | Backlog / nice to have |
| `epic:openapi-gen` | Part of the OpenAPI generation feature |