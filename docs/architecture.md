# Architecture

![Architecture](work_tracker_architecture.png)

Worklog is a modular monolith: one Go service with clear internal modules.

## Components
- **cli/**: Go CLI used on the work PC to log entries
- **server/**: Go API deployed on Cloud Run
  - entries: capture
  - tasks: lifecycle and blockers
  - digest: end-of-day summary and tomorrow's plan
  - notifier: push via FCM
  - jobs: handlers called by the scheduled triggers
- **Neon Postgres**: all data
- **GitHub Actions cron**: calls the jobs endpoints each evening and morning
- **mobile/**: Flutter app to view digests and receive notifications

## Daily flow

```mermaid
sequenceDiagram
    participant CLI as wl CLI
    participant Cron as GitHub cron
    participant API as Go API
    participant DB as Postgres
    participant FCM
    participant App as Flutter app
    CLI->>API: log entry
    API->>DB: store task and event
    Cron->>API: evening job (6 PM Nepal)
    API->>DB: read the day's state
    API->>DB: store digest
    Cron->>API: morning job (8 AM Nepal)
    API->>FCM: send stored plan
    FCM->>App: notification
    App->>API: read digest and tasks
```

## Key decisions
See [decisions](decisions/).
