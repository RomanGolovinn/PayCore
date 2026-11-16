package ledger

type AccountRepository interface {
	AddAccount(account *Account) error
	IsExistAccount(account *Account) (bool, error)
}

type TransactionRepository interface {
	AddTransaction(transaction *Transaction) error
	AddEntry(entry *Entry) error
}
