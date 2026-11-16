package ledger

import (
	"context"
	"errors"
)

type Ledger struct {
	AccountRepository     AccountRepository
	TransactionRepository TransactionRepository
}

func NewLedger(accountRepository AccountRepository, transactionRepository TransactionRepository) *Ledger {
	return &Ledger{AccountRepository: accountRepository, TransactionRepository: transactionRepository}
}

func (l *Ledger) PostTransaction(ctx context.Context, ptc PostTransactionCommand) error {
	var totalSum Amount
	var currency Currency
	entries := make([]Entry, 0, len(ptc.Postings))
	for _, posting := range ptc.Postings {
		if isExist, err := l.AccountRepository.IsAccountExistByID(posting.AccountID); err != nil {
			return errors.New("invalid accountID")
		} else if !isExist {
			return errors.New("invalid accountID")
		}

		entries = append(entries, Entry{posting.AccountID, posting.Money})
		totalSum += posting.Money.Amount
		if currency == "" {
			currency = posting.Money.Currency
		}
		if currency != posting.Money.Currency {
			return errors.New("invalid currency")
		}
	}
	if totalSum != 0 {
		return errors.New("invalid transaction amount")
	}

	transaction := Transaction{
		entries: entries,
	}
	_, err := l.TransactionRepository.AddTransaction(&transaction)
	if err != nil {
		return errors.New("failed to add transaction")
	}

	return nil
}

type PostTransactionCommand struct {
	Postings []Posting
}

type Posting struct {
	AccountID AccountID
	Money     Money
}
