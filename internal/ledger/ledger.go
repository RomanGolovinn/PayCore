package ledger

import "context"

type Ledger struct {
	AccountRepository     AccountRepository
	TransactionRepository TransactionRepository
}

func NewLedger(accountRepository AccountRepository, transactionRepository TransactionRepository) *Ledger {
	return &Ledger{}
}

func (l *Ledger) PostTransaction(ctx context.Context, ptc PostTransactionCommand) error {
	return nil
}

type PostTransactionCommand struct {
	Postings []Posting
}

type Posting struct {
	AccountID AccountID
	Money     Money
}
