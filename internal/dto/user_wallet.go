package dto

import "errors"

type Wallet struct {
	ID      uint64 `form:"id" json:"id"`
	Balance int64  `form:"balance" json:"balance"`
}

type WalletID struct {
	ID uint64 `form:"id" json:"id"`
}

type TransactionID struct {
	ID string `form:"id" json:"id"`
}

type Transaction struct {
	ID     string  `form:"transactionId" json:"transactionId"`
	State  string  `form:"state" json:"state"`
	Amount float64 `form:"amount" json:"amount"`
}

func (t *Transaction) Validate() error {

	if t.ID == "" {
		return errors.New("transactionId not valid")
	}

	if t.Amount <= 0 {
		return errors.New("amount not valid")
	}

	switch t.State {
	case "win", "lost":
		break
	default:
		return errors.New("state not valid")
	}

	return nil
}
