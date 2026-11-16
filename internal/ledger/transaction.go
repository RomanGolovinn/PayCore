package ledger

import "uuid"

type Transaction struct {
	TransactionID TransactionID
	entries       []Entry
}

type TransactionID uuid.UUID
