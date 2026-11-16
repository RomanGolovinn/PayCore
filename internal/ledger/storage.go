package ledger

type AccountRepository interface {
	AddAccount(account *Account) error
	IsAccountExist(account *Account) (bool, error)
	IsAccountExistByID(accountID AccountID) (bool, error)
}

type TransactionRepository interface {
	AddTransaction(transaction *Transaction) (TransactionID, error)
	AddEntry(entry *Entry) error
}
