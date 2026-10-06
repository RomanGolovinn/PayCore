# PayCore

PayCore is a payment platform written in Go.

The project is a long-term playground for designing and building financial infrastructure and for studying backend and distributed-systems engineering through real problems.

## What it explores

* payment processing
* double-entry accounting
* financial ledgers
* idempotency and concurrency
* distributed systems
* event-driven architecture
* fault tolerance and reliability
* security
* observability
* horizontal scaling
* reconciliation and settlement
* multi-currency payments
* risk management
* chargebacks and disputes

The goal is not a simplified demo, but a gradual move toward production-grade architecture and engineering practices.

## Architecture

The initial architecture is a **modular monolith**.

The system is designed around explicit domain boundaries, so individual components can be extracted into independent services later — when there is a real architectural reason to do so.

Architecture documentation lives in [`docs/`](docs/).

## Engineering Principles

* PostgreSQL is the source of truth for financial data.
* Financial operations must be auditable.
* Money is represented using integer minor units, never floating-point.
* Financial operations must be idempotent.
* Domain invariants are enforced explicitly.
* Distributed failures are treated as normal cases.
* Infrastructure complexity is introduced only when it solves a real problem.
* Security-critical functionality relies on established standards and libraries, not custom cryptography.
* Significant architectural decisions are documented.

## Development Status

🚧 **Early development**

The project is currently in the architecture and domain-design stage.

The initial implementation has not started yet. The system will be built incrementally, starting with the core financial domain and evolving toward more complex payment-processing and distributed-systems capabilities.

## Repository Structure

```text
paycore/
├── .github/        # GitHub configuration
├── cmd/            # application entrypoints
├── deployments/    # deployment and infrastructure configuration
├── docs/           # architecture and engineering documentation
├── internal/       # application and domain code
├── migrations/     # database migrations
├── proto/          # Protocol Buffer definitions
├── README.md
├── CONTRIBUTING.md
└── .gitignore
```

The repository structure will evolve as the system grows.

## Development Workflow

Development rules, branch naming conventions, commit conventions, and Pull Request guidelines are documented in [`CONTRIBUTING.md`](CONTRIBUTING.md).

## License

TBD
