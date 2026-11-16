package ledger

type Entry struct {
	AccountID AccountID
	Money     Money
}

type Money struct {
	Amount   Amount
	Currency Currency
}

type Amount int64

func NewEntry(accountID AccountID, money Money) *Entry {
	return &Entry{AccountID: accountID, Money: money}
}
