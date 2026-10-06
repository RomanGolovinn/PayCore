# Domain Boundaries

This document describes the initial domain boundaries of PayCore.

The boundaries are intentionally coarse at this stage. They describe the responsibilities of the major business domains without prescribing a microservice architecture.

The current architecture is a **modular monolith**. A domain may become an independent service in the future if there is a concrete architectural reason to extract it.

## Initial Domains

PayCore currently consists of five primary domains:

1. Identity
2. Payments
3. Risk
4. Ledger
5. Settlement

```text
                         PayCore
                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
    Identity            Payments              Risk
                            │
                          Ledger
                            │
                        Settlement
```

The diagram represents logical ownership and responsibilities, not necessarily direct runtime dependencies.

---

## Identity

Identity is responsible for the participants of the system and their access to it.

### Responsibilities

* customers
* merchants
* users
* authentication
* API credentials
* roles and permissions

Identity answers questions such as:

* Who is this actor?
* Which merchant does this API credential belong to?
* Is this actor allowed to perform this operation?

Identity does not own payment or financial state.

---

## Payments

Payments manages payment intents and their lifecycle.

### Responsibilities

* payment creation
* payer and merchant association
* amount and currency
* payment methods
* payment attempts
* authorization
* capture
* cancellation
* refunds
* payment state transitions

Payments answers the question:

> What is happening with this payment?

Payments does not own the financial ledger and does not determine the final accounting representation of a financial operation.

---

## Risk

Risk evaluates whether a financial operation should be allowed, rejected, or subjected to additional checks.

### Responsibilities

* risk rules
* risk scoring
* velocity limits
* suspicious activity detection
* payment risk decisions

Risk answers the question:

> Is this operation acceptable from a risk perspective?

Risk does not own authentication or access control. It also does not directly modify the financial ledger.

---

## Ledger

Ledger is responsible for the financial accounting of the system.

It is not an application log. It is the source of truth for the financial history of accounts.

### Responsibilities

* accounts
* financial transactions
* ledger entries
* double-entry accounting
* financial balances derived from ledger entries
* immutable financial history

For a transfer of value, the ledger records both sides of the operation.

For example:

```text
Customer account    -100 USD
Merchant account    +100 USD
-----------------------------
                     0 USD
```

The core invariant is that a balanced financial transaction does not create or destroy value:

```text
sum(entries) = 0
```

Ledger does not own the business lifecycle of a payment.

It records the financial consequences of operations initiated by other domains.

---

## Settlement

Settlement is responsible for the process of settling funds between participants and external financial systems.

A payment being successfully processed does not necessarily mean that the merchant has already received settled funds.

### Responsibilities

* settlement batches
* merchant settlements
* payouts
* settlement schedules
* reconciliation with external settlement records

Settlement answers the question:

> When and how are funds settled to the merchant or another financial participant?

Settlement is separate from Payments because payment processing and settlement may happen at different times.

---

## Domain Ownership

Each domain should own its own business rules and state.

Other domains should not directly modify another domain's internal state.

For example:

```text
Payments
    │
    │ payment succeeded
    ▼
Ledger
    │
    │ financial records
    ▼
Settlement
```

The exact communication mechanism between domains is intentionally not defined yet.

It may initially be implemented through direct application-level calls inside the modular monolith. Later, some interactions may become asynchronous events or service-to-service communication.

---

## Current Scope

These five domains are the initial MVP boundaries.

The model is expected to evolve as the system grows.

Potential future domains or components include:

* Payment Provider Integration
* Payment Orchestration
* Webhooks
* Fraud Detection
* Foreign Exchange
* Disputes
* Chargebacks
* Notifications
* Compliance
* Reporting

These are deliberately not treated as separate domains yet.

A new domain should be introduced when its responsibilities, invariants, or lifecycle become sufficiently independent to justify a separate boundary.

## Architectural Principle

The goal is not to predict the final architecture in advance.

The goal is to establish clear boundaries that allow the architecture to evolve without turning individual domains into tightly coupled components or global "god objects".
