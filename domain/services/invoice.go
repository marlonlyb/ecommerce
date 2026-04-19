package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/mlbautomation/Ecommmerce_MLB/domain/ports/invoice"
	"github.com/mlbautomation/Ecommmerce_MLB/domain/ports/purchaseorder"
	"github.com/mlbautomation/Ecommmerce_MLB/model"
)

type Invoice struct {
	Repository              invoice.Repository
	RepositoryInvoiceReport invoice.RepositoryInvoiceReport
	/* services.Invoice usa un puerto creado en purchaseorder
	del tipo Service específico para Invoice que solo
	tiene el Validate(m *model.PurchaseOrder) error */
	ServicePurchaseorder purchaseorder.ServiceInvoice
}

func NewInvoice(ir invoice.Repository, irir invoice.RepositoryInvoiceReport, posi purchaseorder.ServiceInvoice) *Invoice {
	return &Invoice{Repository: ir, RepositoryInvoiceReport: irir, ServicePurchaseorder: posi}
}

func (i *Invoice) Create(po *model.PurchaseOrder) error {

	err := i.ServicePurchaseorder.Validate(po)
	if err != nil {
		return fmt.Errorf("i.ServiceInvoiceProduct.Validate(po): %w", err)
	}

	invoice, invoiceDetails, err := invoiceFromPurchaseOrder(po)
	if err != nil {
		return fmt.Errorf("%s %w", "invoiceFromPurchaseOrder(po)", err)
	}

	//que pasa si cuando estoy insertando invooce details genero un error?
	//debo trabajar con transacciones en la base de datos.
	err = i.Repository.Create(&invoice, invoiceDetails)
	if err != nil {
		return fmt.Errorf("%s %w", "i.Repository.Create(&invoice, invoiceDetails)", err)
	}

	return nil
}

func (i Invoice) GetByUserID(userID uuid.UUID) (model.InvoicesReport, error) {
	invoicesHead, err := i.RepositoryInvoiceReport.HeadsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("invoice: %w", err)
	}

	var invoicesReport model.InvoicesReport
	for _, invoiceHead := range invoicesHead {
		invoiceDetails, err := i.RepositoryInvoiceReport.AllDetailsByInvoiceID(invoiceHead.Invoice.ID)
		if err != nil {
			return nil, fmt.Errorf("%s %w", "i.RepositoryInvoiceReport.AllDetailsByInvoiceID(invoiceHead.Invoice.ID)", err)
		}

		invoiceHead.InvoiceDetailsReport = invoiceDetails
		invoicesReport = append(invoicesReport, invoiceHead)
	}

	return invoicesReport, nil
}

// GetAll returns a model.Invoices according to filters and sorts
func (i Invoice) GetAll() (model.InvoicesReport, error) {
	invoices, err := i.RepositoryInvoiceReport.AllHead()
	if err != nil {
		return nil, fmt.Errorf("invoice: %w", err)
	}

	var invoicesReport model.InvoicesReport
	for _, v := range invoices {
		invoiceDetails, err := i.RepositoryInvoiceReport.AllDetailsByInvoiceID(v.Invoice.ID)
		if err != nil {
			return nil, fmt.Errorf("%s %w", "i.RepositoryInvoiceReport.AllDetailsByInvoiceID(v.Invoice.ID)", err)
		}

		v.InvoiceDetailsReport = invoiceDetails
		invoicesReport = append(invoicesReport, v)
	}

	return invoicesReport, nil
}

func invoiceFromPurchaseOrder(po *model.PurchaseOrder) (model.Invoice, model.InvoiceDetails, error) {
	ID, err := uuid.NewUUID()
	if err != nil {
		return model.Invoice{}, nil, fmt.Errorf("%s %w", "uuid.NewUUID()", err)
	}

	invoice := model.Invoice{
		ID:              ID,
		UserID:          po.UserID,
		PurchaseOrderID: po.ID,
		CreatedAt:       time.Now().Unix(),
	}

	var products model.ProductsToPurchases
	err = json.Unmarshal(po.Products, &products)
	if err != nil {
		return model.Invoice{}, nil, fmt.Errorf("%s %w", "json.Unmarshal()", err)
	}

	var invoiceDetails model.InvoiceDetails
	for _, v := range products {
		detailID, err := uuid.NewUUID()
		if err != nil {
			return model.Invoice{}, nil, fmt.Errorf("%s %w", "uuid.NewUUID()", err)
		}

		// importante validar que el precio sea el que está hoy activo
		detail := model.InvoiceDetail{
			ID:        detailID,
			InvoiceID: invoice.ID,
			ProductID: v.ProductID,
			Amount:    v.Amount,
			UnitPrice: v.UnitPrice,
			CreatedAt: time.Now().Unix(),
		}

		invoiceDetails = append(invoiceDetails, detail)
	}

	return invoice, invoiceDetails, nil
}
