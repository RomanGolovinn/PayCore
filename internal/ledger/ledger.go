package ledger

import (
	"context"
	"fmt"
	"math"
)

type Ledger struct {
	AccountRepository     AccountRepository
	TransactionRepository TransactionRepository
}

func NewLedger(accountRepository AccountRepository, transactionRepository TransactionRepository) *Ledger {
	return &Ledger{AccountRepository: accountRepository, TransactionRepository: transactionRepository}
}

func (l *Ledger) PostTransaction(ctx context.Context, ptc PostTransactionCommand) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if len(ptc.Postings) < 2 {
		return InvalidPostingsCountError
	}

	var totalSum Amount
	var currency Currency
	entries := make([]Entry, 0, len(ptc.Postings))
	for _, posting := range ptc.Postings {
		if isExist, err := l.AccountRepository.IsAccountExistByID(posting.AccountID); err != nil {
			return fmt.Errorf("check account existence: %w", err)
		} else if !isExist {
			return InvalidAccountError
		}

		if posting.Money.Amount == 0 {
			return InvalidEntryAmountError
		}

		var err error
		currency, err = l.currencyCheck(currency, posting.Money.Currency)
		if err != nil {
			return err
		}
		totalSum, err = addAmount(totalSum, posting.Money.Amount)
		if err != nil {
			return err
		}

		entries = append(entries, Entry{
			AccountID: posting.AccountID,
			Money:     posting.Money,
		})
	}
	if totalSum != 0 {
		return InvalidTransactionAmountError
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	transaction := Transaction{
		entries: entries,
	}
	_, err := l.TransactionRepository.AddTransaction(&transaction)
	if err != nil {
		return fmt.Errorf("add transaction: %w", err)
	}

	return nil
}

func (l *Ledger) currencyCheck(transactionCurrency, entryCurrency Currency) (Currency, error) {
	if !isSupportedCurrency(entryCurrency) {
		return "", InvalidCurrencyError
	}
	if transactionCurrency == "" {
		return entryCurrency, nil
	}
	if entryCurrency != transactionCurrency {
		return "", InvalidCurrencyError
	}
	return transactionCurrency, nil
}

func isSupportedCurrency(currency Currency) bool {
	switch currency {
	case USD, EUR, GBP, JPY, RUB:
		return true
	default:
		return false
	}
}

func addAmount(left, right Amount) (Amount, error) {
	if right > 0 && left > Amount(math.MaxInt64)-right {
		return 0, AmountOverflowError
	}
	if right < 0 && left < Amount(math.MinInt64)-right {
		return 0, AmountOverflowError
	}

	return left + right, nil
}

type PostTransactionCommand struct {
	Postings []Posting
}

type Posting struct {
	AccountID AccountID
	Money     Money
}
