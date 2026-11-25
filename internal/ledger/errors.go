package ledger

import "errors"

var InvalidCurrencyError = errors.New("invalid currency")
var InvalidEntryAmountError = errors.New("invalid entry amount")
var InvalidTransactionAmountError = errors.New("invalid transaction amount")
var InvalidAccountError = errors.New("invalid account")
var InvalidPostingsCountError = errors.New("transaction must contain at least two postings")
var AmountOverflowError = errors.New("amount overflow")
