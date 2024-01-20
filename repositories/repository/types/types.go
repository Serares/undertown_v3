package types

import (
	"errors"
)

var (
	ErrorNotFound          = errors.New("row doesn't exist")
	ErrorAccessingDatabase = errors.New("error while trying to query the database")
)

type TransactionType int

const (
	Sell TransactionType = iota
	Rent
	Default // using this value to call the ToInt method
)

// Will return SELL or RENT strings
func (t TransactionType) String() string {
	return [...]string{"SELL", "RENT"}[t]
}

// Not sure if this is needed
func (t TransactionType) ToInt(stringType string) TransactionType {
	switch stringType {
	case "SELL":
		return Sell
	case "RENT":
		return Rent
	default:
		return Sell
	}
}
