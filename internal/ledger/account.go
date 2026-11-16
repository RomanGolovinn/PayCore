package ledger

import "uuid"

type Account struct {
	AccountID AccountID
	Currency  Currency
}
type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
	JPY Currency = "JPY"
	RUB Currency = "RUB"
)

type AccountID uuid.UUID
