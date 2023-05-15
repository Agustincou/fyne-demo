package repository

import (
	"errors"
	"time"
)

var (
	errUpdateFailed = errors.New("update failed")
	errDeleteFailed = errors.New("delete failed")
)

type Repository interface {
	Migrate() error //To create or load an existent database
	InsertHolding(h Holding) (*Holding, error)
	AllHoldings() ([]Holding, error)
	GetHoldingByID(id int) (*Holding, error)
	UpdateHolding(id int64, updated Holding) error
	DeleteHolding(id int64) error
}

type Holding struct {
	ID            int64     `json:"id"`
	Amount        int       `json:"amount"`
	PurchaseDate  time.Time `json:"purchase_date"`
	PurchasePrice int       `json:"purchase_price"`
}
