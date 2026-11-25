package ledger

import (
	"context"
	"errors"
	"math"
	"testing"
	"uuid"
)

func TestLedgerPostTransaction(t *testing.T) {
	t.Parallel()

	accountA := AccountID(uuid.New())
	accountB := AccountID(uuid.New())
	accountC := AccountID(uuid.New())
	repositoryError := errors.New("repository unavailable")

	tests := []struct {
		name              string
		postings          []Posting
		accountRepository *fakeAccountRepository
		wantErr           error
		wantSaved         bool
		wantEntries       int
	}{
		{
			name: "posts balanced transaction with two entries",
			postings: []Posting{
				posting(accountA, -100, USD),
				posting(accountB, 100, USD),
			},
			accountRepository: existingAccounts(accountA, accountB),
			wantSaved:         true,
			wantEntries:       2,
		},
		{
			name: "posts balanced transaction with multiple entries",
			postings: []Posting{
				posting(accountA, -100, USD),
				posting(accountB, 97, USD),
				posting(accountC, 3, USD),
			},
			accountRepository: existingAccounts(accountA, accountB, accountC),
			wantSaved:         true,
			wantEntries:       3,
		},
		{
			name:              "rejects empty transaction",
			accountRepository: existingAccounts(),
			wantErr:           InvalidPostingsCountError,
		},
		{
			name: "rejects transaction with one posting",
			postings: []Posting{
				posting(accountA, 100, USD),
			},
			accountRepository: existingAccounts(accountA),
			wantErr:           InvalidPostingsCountError,
		},
		{
			name: "rejects missing account",
			postings: []Posting{
				posting(accountA, -100, USD),
				posting(accountB, 100, USD),
			},
			accountRepository: existingAccounts(accountA),
			wantErr:           InvalidAccountError,
		},
		{
			name: "preserves account repository error",
			postings: []Posting{
				posting(accountA, -100, USD),
				posting(accountB, 100, USD),
			},
			accountRepository: &fakeAccountRepository{err: repositoryError},
			wantErr:           repositoryError,
		},
		{
			name: "rejects zero amount",
			postings: []Posting{
				posting(accountA, 0, USD),
				posting(accountB, 0, USD),
			},
			accountRepository: existingAccounts(accountA, accountB),
			wantErr:           InvalidEntryAmountError,
		},
		{
			name: "rejects empty currency",
			postings: []Posting{
				posting(accountA, -100, ""),
				posting(accountB, 100, ""),
			},
			accountRepository: existingAccounts(accountA, accountB),
			wantErr:           InvalidCurrencyError,
		},
		{
			name: "rejects unsupported currency",
			postings: []Posting{
				posting(accountA, -100, Currency("BTC")),
				posting(accountB, 100, Currency("BTC")),
			},
			accountRepository: existingAccounts(accountA, accountB),
			wantErr:           InvalidCurrencyError,
		},
		{
			name: "rejects mixed currencies",
			postings: []Posting{
				posting(accountA, -100, USD),
				posting(accountB, 100, EUR),
			},
			accountRepository: existingAccounts(accountA, accountB),
			wantErr:           InvalidCurrencyError,
		},
		{
			name: "rejects unbalanced transaction",
			postings: []Posting{
				posting(accountA, -100, USD),
				posting(accountB, 90, USD),
			},
			accountRepository: existingAccounts(accountA, accountB),
			wantErr:           InvalidTransactionAmountError,
		},
		{
			name: "rejects amount overflow",
			postings: []Posting{
				posting(accountA, Amount(math.MaxInt64), USD),
				posting(accountB, 1, USD),
				posting(accountC, Amount(math.MinInt64), USD),
			},
			accountRepository: existingAccounts(accountA, accountB, accountC),
			wantErr:           AmountOverflowError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			transactions := &fakeTransactionRepository{}
			ledger := NewLedger(test.accountRepository, transactions)

			err := ledger.PostTransaction(context.Background(), PostTransactionCommand{
				Postings: test.postings,
			})

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("PostTransaction() error = %v, want %v", err, test.wantErr)
			}
			if transactions.called != test.wantSaved {
				t.Fatalf("AddTransaction() called = %v, want %v", transactions.called, test.wantSaved)
			}
			if test.wantSaved && len(transactions.transaction.entries) != test.wantEntries {
				t.Fatalf("saved entries = %d, want %d", len(transactions.transaction.entries), test.wantEntries)
			}
		})
	}
}

func TestLedgerPostTransactionContextCancellation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	transactions := &fakeTransactionRepository{}
	ledger := NewLedger(existingAccounts(), transactions)

	err := ledger.PostTransaction(ctx, PostTransactionCommand{})

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("PostTransaction() error = %v, want %v", err, context.Canceled)
	}
	if transactions.called {
		t.Fatal("AddTransaction() was called for canceled context")
	}
}

func TestLedgerPostTransactionPreservesStorageError(t *testing.T) {
	t.Parallel()

	accountA := AccountID(uuid.New())
	accountB := AccountID(uuid.New())
	storageError := errors.New("write failed")
	transactions := &fakeTransactionRepository{err: storageError}
	ledger := NewLedger(existingAccounts(accountA, accountB), transactions)

	err := ledger.PostTransaction(context.Background(), PostTransactionCommand{
		Postings: []Posting{
			posting(accountA, -100, USD),
			posting(accountB, 100, USD),
		},
	})

	if !errors.Is(err, storageError) {
		t.Fatalf("PostTransaction() error = %v, want wrapped %v", err, storageError)
	}
}

func posting(accountID AccountID, amount Amount, currency Currency) Posting {
	return Posting{
		AccountID: accountID,
		Money: Money{
			Amount:   amount,
			Currency: currency,
		},
	}
}

func existingAccounts(accountIDs ...AccountID) *fakeAccountRepository {
	accounts := make(map[AccountID]bool, len(accountIDs))
	for _, accountID := range accountIDs {
		accounts[accountID] = true
	}
	return &fakeAccountRepository{accounts: accounts}
}

type fakeAccountRepository struct {
	accounts map[AccountID]bool
	err      error
}

func (r *fakeAccountRepository) AddAccount(*Account) error {
	return nil
}

func (r *fakeAccountRepository) IsAccountExist(account *Account) (bool, error) {
	return r.IsAccountExistByID(account.AccountID)
}

func (r *fakeAccountRepository) IsAccountExistByID(accountID AccountID) (bool, error) {
	if r.err != nil {
		return false, r.err
	}
	return r.accounts[accountID], nil
}

type fakeTransactionRepository struct {
	transaction Transaction
	called      bool
	err         error
}

func (r *fakeTransactionRepository) AddTransaction(transaction *Transaction) (TransactionID, error) {
	r.called = true
	if r.err != nil {
		return TransactionID{}, r.err
	}
	r.transaction = *transaction
	return TransactionID(uuid.New()), nil
}

func (r *fakeTransactionRepository) AddEntry(*Entry) error {
	return nil
}
