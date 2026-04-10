package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

type PurchaseOrder struct {
	ID        uuid.UUID       `json:"id"`
	UserID    uuid.UUID       `json:"user_id"`
	Products  json.RawMessage `json:"products"`
	CreatedAt int64           `json:"created_at"`
	UpdatedAt int64           `json:"updated_at"`
}

func (p PurchaseOrder) HasID() bool {
	return p.ID != uuid.Nil
}

type PurchaseOrders []PurchaseOrder

func (p PurchaseOrders) IsEmpty() bool {
	return len(p) == 0
}

// Ésta es la estructura que recibo en el body
type ProductsToPurchase struct {
	ProductID uuid.UUID `json:"product_id"`
	Amount    uint      `json:"amount"`
	UnitPrice float64   `json:"unit_price"`
}

type ProductsToPurchases []ProductsToPurchase
