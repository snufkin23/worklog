# Architecture

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
- **GitHub Actions cron**: calls the jobs endpoints at 6 PM and 8 AM
- **mobile/**: Flutter app to view digests and receive notifications

## Flows
1. Capture: CLI -> API -> entries/tasks -> database
2. Automation: cron -> jobs -> digest (evening) / notifier (morning) -> FCM -> phone
3. Reading: Flutter app -> API -> database
