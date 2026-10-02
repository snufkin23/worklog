# Data model

| Table | Holds |
|---|---|
| tasks | id (short number), title, status (open, blocked, done, dropped), blocker reason, important flag, created, updated, closed |
| events | id, optional task id, type (created, progress, blocked, unblocked, done, dropped, note), text, created |
| digests | date (unique), summary, plan, generated at, sent at |
| devices | push token, last seen |

## Rules
- Events are append-only. History is never edited.
- A note is an event with no task.
- Timestamps are stored in UTC; a "day" is a Nepal day (UTC+5:45).
- One digest per date. `sent_at` is set once, so retried jobs do nothing twice.
