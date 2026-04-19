package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/mlbautomation/Ecommmerce_MLB/domain/ports/product"
	"github.com/mlbautomation/Ecommmerce_MLB/domain/ports/purchaseorder"
	"github.com/mlbautomation/Ecommmerce_MLB/model"
)

type PurchaseOrder struct {
	Repository purchaseorder.Repository
	/* services.PurchaseOrder usa un puerto creado en product
	del tipo Repository específico para PurchaseOrder que solo
	tiene el GetByID(ID uuid.UUID) (model.Product, error) */
	RepositoryProduct product.RepositoryPurchaseOrder
}

func NewPurchaseOrder(por purchaseorder.Repository, prpo product.RepositoryPurchaseOrder) *PurchaseOrder {
	return &PurchaseOrder{Repository: por, RepositoryProduct: prpo}
}

func (p *PurchaseOrder) Create(m *model.PurchaseOrder) error {

	err := p.Validate(m)
	if err != nil {
		return fmt.Errorf("Validate(): %w", err)
	}

	ID, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("%s %w", "uuid.NewUUID()", err)
	}
	m.ID = ID

	m.CreatedAt = time.Now().Unix()

	err = p.Repository.Create(m)
	if err != nil {
		return fmt.Errorf("%s %w", "Repository.Create(m)", err)
	}

	return nil
}

func (p *PurchaseOrder) GetByID(ID uuid.UUID) (model.PurchaseOrder, error) {
	purchaseOrder, err := p.Repository.GetByID(ID)
	if err != nil {
		return model.PurchaseOrder{}, fmt.Errorf("%s %w", "Repository.GetByID(ID)", err)
	}

	return purchaseOrder, nil
}

func (p *PurchaseOrder) Validate(m *model.PurchaseOrder) error {

	if len(m.Products) == 0 {
		return errors.New("the list of products can't be empty")
	}

	var ptps []model.ProductsToPurchase
	err := json.Unmarshal(m.Products, &ptps)
	if err != nil {
		return fmt.Errorf("%s %w", "json.Unmarshal()", err)
	}

	for i, v := range ptps {
		if v.ProductID == uuid.Nil {
			return errors.New("the product id can't be empty")
		}

		mp, err := p.RepositoryProduct.GetByID(v.ProductID)
		if err != nil {
			return fmt.Errorf("%s %w", "there is an ProductID invalid", err)
		}
		ptps[i].UnitPrice = mp.Price

		if v.Amount < 1 {
			return errors.New("the amount of products can't ve less tha 1")
		}
	}

	m.Products, err = json.Marshal(ptps)
	if err != nil {
		return fmt.Errorf("%s %w", "json.Marshal()", err)
	}

	return nil
}

// Aquí ya no verifico nada, se supono que se verificó antes de crear el PurchaseOrder
func (p *PurchaseOrder) TotalAmount(m model.PurchaseOrder) float64 {

	var ptps []model.ProductsToPurchase
	err := json.Unmarshal(m.Products, &ptps)
	if err != nil {
		return 0
	}

	var total float64

	for _, v := range ptps {

		mp, err := p.RepositoryProduct.GetByID(v.ProductID)
		if err != nil {
			return 0
		}

		total += float64(v.Amount) * mp.Price
	}

	return total
}
