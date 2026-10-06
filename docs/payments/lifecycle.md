# Payment Lifecycle

This document describes the initial lifecycle of a payment in PayCore.

The lifecycle represents the business state of a payment. It does not define the implementation details of the underlying services or infrastructure.

## Overview

A payment represents an attempt to transfer value from a customer to a merchant.

At a high level:

```text id="p5p6k5"
Customer
    │
    │ create payment
    ▼
Payment Intent
    │
    │ risk evaluation
    ▼
Risk Check
    │
    │ approved
    ▼
Payment Processing
    │
    │ successful
    ▼
Ledger
    │
    │ funds become available for settlement
    ▼
Settlement
    │
    ▼
Merchant
```

The exact interaction between domains will be defined separately.

---

## Payment Intent

A `PaymentIntent` represents the intention to make a payment.

It contains the information required to describe the desired payment, such as:

* customer
* merchant
* amount
* currency
* payment method
* current payment status

A PaymentIntent is not necessarily a completed payment.

For example, a customer may create a payment and then fail during authorization.

---

## Payment States

The initial state machine is:

```text id="o2f6l3"
                    ┌───────────┐
                    │  CREATED  │
                    └─────┬─────┘
                          │
                          ▼
                    ┌───────────┐
                    │ PROCESSING│
                    └─────┬─────┘
                          │
                    ┌─────┴─────┐
                    ▼           ▼
              ┌───────────┐ ┌───────────┐
              │ SUCCEEDED │ │  FAILED   │
              └─────┬─────┘ └───────────┘
                    │
                    ▼
              ┌───────────┐
              │ CAPTURED  │
              └─────┬─────┘
                    │
                    ▼
              ┌───────────┐
              │ REFUNDED  │
              └───────────┘
```

This state machine is intentionally simplified for the initial MVP.

The set of states and transitions may evolve as payment methods and provider integrations are added.

---

## Creating a Payment

The customer or merchant creates a PaymentIntent.

Example:

```text id="f2pr6b"
Customer:  customer-123
Merchant:  merchant-456
Amount:    100 USD
```

The PaymentIntent initially enters the `CREATED` state.

At this stage, no successful payment has occurred.

---

## Risk Evaluation

Before processing the payment, PayCore may ask the Risk domain to evaluate the operation.

```text id="4y4xuw"
Payment
   │
   ▼
Risk
   │
   ├── APPROVE
   ├── REJECT
   └── REVIEW
```

A rejected payment does not proceed to financial processing.

Risk does not change the payment state directly. The Payments domain remains responsible for the payment lifecycle.

---

## Processing

Once the payment is allowed to proceed, it enters `PROCESSING`.

At this stage PayCore may communicate with an external payment provider.

The provider may:

* approve the operation;
* reject the operation;
* fail to respond;
* return a temporary error;
* require additional processing.

The payment remains `PROCESSING` while its final result is unknown.

---

## Successful Payment

When the payment is successfully processed, it becomes eligible for financial accounting.

The Payments domain records the successful business operation.

The Ledger then records its financial consequences.

For a simplified $100 payment:

```text id="2kgm8s"
Customer account    -100 USD
Merchant account    +100 USD
-----------------------------
                     0 USD
```

The Ledger is responsible for the financial representation of the operation.

Payments remains responsible for the payment's business state.

---

## Failed Payment

A payment may fail during risk evaluation or external processing.

Examples include:

* risk rejection;
* insufficient funds;
* provider rejection;
* expired authorization;
* permanent provider error.

A failed payment does not create a successful financial transfer.

```text id="5q6qf4"
PROCESSING
     │
     ▼
  FAILED
```

The exact failure reason should be recorded separately from the high-level payment state.

---

## Capture and Refund

The initial MVP may simplify authorization and capture into a single flow.

Later, the payment lifecycle can be extended to support explicit authorization and capture:

```text id="2j4m9f"
CREATED
   ↓
AUTHORIZED
   ↓
CAPTURED
```

A captured payment may later be refunded:

```text id="x1u7as"
CAPTURED
   ↓
REFUNDED
```

Partial refunds will also be supported by the domain model in the future.

For example:

```text id="7p6c8h"
Payment: 100 USD

Refund:   30 USD

Remaining refundable amount: 70 USD
```

---

## Settlement

A successful payment and a completed settlement are separate concepts.

After a payment is successfully processed, funds may remain in a pending or unsettled state.

Later, the Settlement domain may include the funds in a settlement batch and make them available to the merchant.

Simplified:

```text id="l5e6as"
Payment
   │
   ▼
Successful
   │
   ▼
Unsettled funds
   │
   │ settlement cycle
   ▼
Settled
   │
   ▼
Merchant payout
```

The exact settlement model depends on the payment method and external financial infrastructure.

---

## Important Invariants

The payment lifecycle must preserve several invariants.

### Invalid state transitions are forbidden

For example:

```text id="7l2f5r"
FAILED → CAPTURED
REFUNDED → CAPTURED
```

must not be possible.

### A payment must not be captured twice

Repeated requests must not create duplicate financial operations.

### A refund cannot exceed the captured amount

For a captured payment of `100 USD`:

```text id="x6a3a8r"
refund(100 USD) → allowed
refund(30 USD)  → allowed
refund(101 USD) → rejected
```

### Financial operations must be idempotent

Retrying the same operation must not result in duplicate financial effects.

---

## Future Extensions

The lifecycle can later be extended with:

* explicit authorization and capture;
* multiple payment attempts;
* payment provider routing;
* asynchronous provider webhooks;
* retries;
* 3-D Secure;
* partial captures;
* partial refunds;
* chargebacks;
* disputes;
* payment method-specific states.

The lifecycle should evolve without breaking the fundamental distinction between:

1. payment business state;
2. financial accounting;
3. settlement.
