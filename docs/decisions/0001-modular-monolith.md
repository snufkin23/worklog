# 0001: Modular monolith instead of microservices

## Context
Single user, single developer, a handful of simple jobs.

## Decision
One deployable Go service with internal modules and clean boundaries.

## Why
Microservices solve team-scale and independent-scaling problems this project
does not have. They would add network calls and operational overhead with no benefit.

## Consequences
A module can be extracted into its own service later if a real need appears.
