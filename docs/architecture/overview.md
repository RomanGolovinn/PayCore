# Architecture Overview

## Architecture Goals

PayCore is designed as a long-term engineering project for studying and implementing production-grade payment infrastructure.

The architecture prioritizes:

* explicit domain boundaries;
* correctness of financial operations;
* reliability;
* auditability;
* idempotency;
* controlled complexity;
* incremental evolution.

The architecture should support increasing system complexity without requiring the entire system to be redesigned from scratch.

## Initial Architecture

The initial implementation is a **modular monolith**.

All domains run inside one application and are developed in one repository.

```text
                     PayCore
                        │
        ┌───────────────┼────────────────┐
        │               │                │
     Identity        Payments          Risk
                        │
                        ↓
                     Ledger
                        │
                        ↓
                   Settlement
```

This is a logical view of the system, not a deployment diagram.

The initial architecture does not require separate microservices for each domain.

## Why a Modular Monolith

A modular monolith provides explicit boundaries without introducing distributed-system complexity before it is necessary.

The goal is to make the internal boundaries strong enough that individual modules can be extracted later if necessary.

## Domain Boundaries

The initial system contains five primary domains:

* Identity
* Payments
* Risk
* Ledger
* Settlement

Their responsibilities are documented in [`domains.md`](domains.md).

Each domain owns its own business rules and state.

Other domains should interact through explicit interfaces rather than modifying another domain's internal state directly.

## Dependency Direction

The initial dependency direction is conceptually:

```text
Transport
    ↓
Application
    ↓
Domain
    ↓
Infrastructure
```

The exact Go package structure may evolve during implementation.

## Domain Interaction

In the modular monolith, domains may initially communicate through direct in-process application calls.

For example:

```text
Payment Application
        │
        ├──→ Risk
        │
        └──→ Ledger
```

This does not mean the domains are allowed to share arbitrary internal structures.

Communication should happen through explicit contracts.

Later, some interactions may move to asynchronous events or service-to-service communication when there is a concrete architectural reason.

The mechanism of communication is therefore deliberately not fixed at this stage.

## Source of Truth

Different domains own different kinds of state.

### Identity

Source of truth for:

* users;
* customers;
* merchants;
* credentials;
* permissions.

### Payments

Source of truth for:

* payment intents;
* payment attempts;
* payment state;
* authorization and capture state.

### Risk

Source of truth for:

* risk evaluations;
* risk decisions;
* risk rules and limits.

### Ledger

Source of truth for:

* financial transactions;
* ledger entries;
* financial history;
* account balances derived from entries.

### Settlement

Source of truth for:

* settlement batches;
* settlement state;
* payout state;
* reconciliation state.

A domain must not silently duplicate another domain's authoritative state and treat the copy as authoritative.

## Data Storage

PostgreSQL is the initial primary database.

Financial data is stored transactionally and must preserve the invariants defined by the Ledger domain.

The initial system does not require a separate database for every domain.

Logical ownership is more important than physical database separation at this stage.

If the system is later split into independent services, data ownership can also be physically separated where appropriate.

## External Systems

Payment providers and other external financial systems are treated as unreliable dependencies.

External operations may:

* timeout;
* fail;
* return temporary errors;
* return duplicate responses;
* complete successfully while the client does not receive the response.

The architecture must therefore be designed around retries, idempotency and reconciliation.

The first implementation can use a fake payment provider for development and testing.

Real provider integrations can be introduced later behind an explicit abstraction.

## Reliability Principles

The architecture assumes that failures are normal.

Examples include:

* application crashes;
* database connection failures;
* network timeouts;
* duplicate requests;
* duplicated messages;
* delayed webhooks;
* external provider failures.

Important financial operations must therefore be:

* atomic where necessary;
* idempotent;
* auditable;
* safe to retry.

The system should prefer explicit recovery mechanisms over assumptions that failures will not happen.

## Events

The initial architecture does not require a distributed event bus.

Events may be introduced when asynchronous processing, integration, or decoupling provides a concrete benefit.

A likely evolution is:

```text
Domain operation
      ↓
Database transaction
      ↓
Transactional Outbox
      ↓
Event Bus
      ↓
Consumers
```

The transactional outbox pattern will be introduced before relying on an external message broker for important domain events.

## Infrastructure Evolution

Infrastructure complexity should grow together with actual system requirements.

A possible evolution is:

```text
Modular Monolith
      ↓
Docker
      ↓
Transactional Outbox
      ↓
Event Bus
      ↓
Horizontal Scaling
      ↓
Service Extraction
      ↓
Kubernetes / Multi-Service Deployment
```

This is an evolutionary path, not a mandatory roadmap.

Each new infrastructure component must solve a concrete problem.

## Security

Security is treated as a cross-cutting architectural concern.

The system will eventually require:

* authentication;
* authorization;
* API credentials;
* TLS;
* secret management;
* key rotation;
* rate limiting;
* audit logging;
* protection against replay and duplicate requests;
* secure handling of sensitive payment data.

Security-critical cryptographic functionality must rely on established standards and libraries rather than custom cryptographic implementations.

## Observability

The system will eventually provide:

* structured logs;
* metrics;
* distributed tracing;
* request correlation;
* business metrics;
* audit information.

Observability should make it possible to understand both technical failures and financial state transitions.

The initial implementation should keep observability simple and add more infrastructure as the system becomes more complex.

## Architectural Evolution

The architecture is intentionally not considered final.

New domains, infrastructure components, or deployment boundaries should be introduced when real requirements justify them.

Examples:

* Payment Provider Integration may become its own module.
* Payment Orchestration may be separated when routing logic becomes complex.
* Webhooks may require asynchronous processing.
* Risk may require independent scaling.
* Ledger may require stronger isolation and specialized infrastructure.
* Service extraction may become appropriate when independent deployment or scaling becomes valuable.

The system should evolve from real constraints rather than from an assumed target architecture.

## Core Architectural Principles

PayCore follows these principles:

1. Keep domain boundaries explicit.
2. Keep financial state authoritative and auditable.
3. Prefer a modular monolith before distributed deployment.
4. Introduce infrastructure only when it solves a real problem.
5. Treat external systems and networks as unreliable.
6. Make important operations idempotent and safe to retry.
7. Keep domain logic independent from transport and infrastructure details.
8. Avoid shared mutable state between domains.
9. Document significant architectural decisions.
10. Evolve the architecture as system requirements become clearer.

The architecture is therefore designed to support complexity without introducing unnecessary complexity in advance.
