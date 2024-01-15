package types

import "errors"

var (
	ErrorNotFound          = errors.New("row doesn't exist")
	ErrorAccessingDatabase = errors.New("error while trying to query the database")
)

type TransactionType int

const (
	Sell TransactionType = iota
	Rent
)

// Will return SELL or RENT strings
func (t TransactionType) String() string {
	return [...]string{"SELL", "RENT"}[t]
}