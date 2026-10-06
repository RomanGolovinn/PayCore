# Ledger Accounting Model

## Purpose

The Ledger is the accounting system of PayCore.

It records the financial consequences of business operations and provides the source of truth for financial history.

The Ledger is not an application log and does not own the lifecycle of payments.

Its responsibility is to answer questions such as:

* how much value was credited or debited;
* which accounts were affected;
* when a financial operation was recorded;
* what the financial history of an account looks like.

## Core Concepts

The initial Ledger model consists of four main concepts:

* **Account** — a balance-holding entity within the financial system.
* **Ledger Transaction** — an atomic financial operation.
* **Ledger Entry** — an individual debit or credit affecting an account.
* **Money** — an integer amount together with a currency.

### Account

An account represents a place where financial value is recorded.

Examples:

* customer account;
* merchant account;
* platform account;
* fee account;
* clearing or settlement account.

The exact account types will evolve as the system grows.

An account belongs to the Ledger domain and its financial state must not be modified directly by other domains.

### Ledger Transaction

A Ledger Transaction represents one atomic financial operation.

Examples:

* payment posting;
* refund;
* fee;
* settlement;
* transfer between accounts.

All entries belonging to one Ledger Transaction are committed atomically.

A transaction is either fully recorded or not recorded at all.

### Ledger Entry

A Ledger Entry represents the effect of a Ledger Transaction on a single account.

A transaction normally contains at least two entries.

For example, a $100 payment:

```text
Customer account    -100 USD
Merchant account    +100 USD
-----------------------------
                       0 USD
```

The entries describe the movement of financial value between accounts.

## Double-Entry Invariant

PayCore uses double-entry accounting.

For every balanced Ledger Transaction:

```text
sum(entries) = 0
```

The invariant applies independently for each currency.

For example:

```text
Customer account    -100 USD
Merchant account    +100 USD
-----------------------------
                       0 USD
```

is valid.

The following is not valid:

```text
Customer account    -100 USD
Merchant account     +90 USD
```

because the transaction does not balance.

Different currencies cannot be used to balance each other:

```text
-100 USD + 100 EUR
```

is not a balanced transaction.

Currency conversion requires an explicit FX operation with its own financial model.

## Money Representation

Money is represented using integer minor units.

For example:

```text
$10.50 = 1050 cents
```

Floating-point numbers are not used for financial amounts.

A monetary value consists of:

```text
amount
currency
```

The amount is an integer and the currency is explicit.

This prevents precision errors and makes financial calculations deterministic.

## Immutability

Recorded financial history must be immutable.

Ledger Entries and completed Ledger Transactions must not be edited or deleted to change historical results.

If an operation needs to be corrected, the system records a new compensating transaction.

For example, instead of modifying:

```text
-100 USD
+100 USD
```

a correction can record:

```text
+100 USD
-100 USD
```

as a separate transaction.

This preserves the complete financial history.

## Balance

An account balance is derived from its Ledger Entries.

Conceptually:

```text
balance = sum(account entries)
```

The system may later introduce cached or materialized balances for performance.

Such values must remain derived data.

The Ledger Entries remain the authoritative financial history.

## Atomicity

All entries belonging to one Ledger Transaction must be persisted atomically.

A partial transaction must never become visible.

For example, this operation:

```text
Customer -100 USD
Merchant +100 USD
```

must either create both entries or create neither.

A state where only one side exists would violate the accounting invariant.

## Idempotency

Financial operations must be idempotent.

Retrying the same operation must not create duplicate financial effects.

For example, if a payment is accidentally submitted to the Ledger twice with the same idempotency key, the Ledger must not record two payments.

Instead, the second request must resolve to the already-recorded financial operation.

Idempotency is especially important because distributed systems can retry operations after timeouts or connection failures.

## Concurrency

The Ledger must remain correct when multiple financial operations affect the same account concurrently.

The implementation must prevent race conditions that could result in:

* incorrect balances;
* duplicated entries;
* lost updates;
* violated accounting invariants.

The exact concurrency-control strategy is an implementation concern and will be defined when the Ledger is implemented.

## Relationship With Payments

The Payments domain owns the payment lifecycle.

The Ledger owns the financial consequences of that lifecycle.

For example:

```text
Payment succeeds
      ↓
Payments determines that money should be moved
      ↓
Ledger records the financial transaction
```

Payments must not directly modify Ledger state.

Likewise, the Ledger must not determine whether a payment should be authorized or captured.

This separation keeps business lifecycle and financial accounting independent.

## Refunds

A refund is represented as a new financial transaction.

The original transaction is not modified.

For example:

```text
Original payment:

Customer    -100 USD
Merchant    +100 USD
```

A full refund creates a compensating transaction:

```text
Customer    +100 USD
Merchant    -100 USD
```

The system must ensure that the total refunded amount cannot exceed the amount eligible for refund.

Partial refunds will be supported later.

## Future Extensions

The initial accounting model is intentionally small.

As PayCore evolves, the Ledger may support:

* platform fees;
* multiple system accounts;
* clearing accounts;
* settlement accounts;
* partial refunds;
* chargebacks;
* disputes;
* multi-currency accounting;
* FX transactions;
* pending and available balances;
* accounting periods;
* reconciliation;
* audit and reporting views.

These features should extend the model without weakening its core invariants.

## Core Invariants

The initial Ledger implementation must preserve the following invariants:

1. Every balanced financial transaction has zero net amount per currency.
2. A financial transaction is committed atomically.
3. Historical Ledger Entries are immutable.
4. Financial operations are idempotent.
5. Money is represented using integer minor units.
6. Currency is explicit and cannot be ignored.
7. Account balances are derived from Ledger history.
8. Other domains cannot directly modify Ledger state.

The Ledger is therefore the financial source of truth of PayCore.
