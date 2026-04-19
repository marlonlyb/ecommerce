package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mlbautomation/Ecommmerce_MLB/model"
)

const (
	psqlHeadInvoiceReport = `
	SELECT invoices.id, invoices.user_id, invoices.purchase_order_id, invoices.created_at,
	invoices.updated_at, users.id, users.email, users.is_admin, users.details, users.created_at,
	users.updated_at
	FROM invoices INNER JOIN users
	ON invoices.user_id = users.id
	`
	psqlDetailInvoiceReport = `
	SELECT invoice_details.id, invoice_details.invoice_id, invoice_details.product_id,
	invoice_details.amount, invoice_details.unit_price, invoice_details.created_at,
	invoice_details.updated_at, products.id, products.product_name, products.price,
	products.images, products.description, products.features, products.created_at,
	products.updated_at
	FROM invoice_details INNER JOIN products ON invoice_details.product_id = products.id
	WHERE invoice_details.invoice_id = $1
	`
)

type InvoiceReport struct {
	db *pgxpool.Pool
}

func NewInvoiceReport(db *pgxpool.Pool) InvoiceReport {
	return InvoiceReport{db: db}
}

// Este no veo donde se utilice...
func (ir InvoiceReport) HeadByInvoiceID(ID uuid.UUID) (model.InvoiceReport, error) {
	row := ir.db.QueryRow(
		context.Background(),
		psqlHeadInvoiceReport+" WHERE invoices.ID = $1",
		ID,
	)
	return ir.scanHead(row)
}

func (ir InvoiceReport) HeadsByUserID(userID uuid.UUID) (model.InvoicesReport, error) {
	rows, err := ir.db.Query(
		context.Background(),
		psqlHeadInvoiceReport+" WHERE users.ID = $1",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resp model.InvoicesReport
	for rows.Next() {
		r, err := ir.scanHead(rows)
		if err != nil {
			return nil, err
		}

		resp = append(resp, r)
	}

	return resp, nil
}

func (ir InvoiceReport) AllHead() (model.InvoicesReport, error) {
	rows, err := ir.db.Query(
		context.Background(),
		psqlHeadInvoiceReport,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resp model.InvoicesReport
	for rows.Next() {
		r, err := ir.scanHead(rows)
		if err != nil {
			return nil, err
		}

		resp = append(resp, r)
	}

	return resp, nil
}

func (ir InvoiceReport) AllDetailsByInvoiceID(ID uuid.UUID) (model.InvoiceDetailsReports, error) {
	rows, err := ir.db.Query(
		context.Background(),
		psqlDetailInvoiceReport,
		ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resp model.InvoiceDetailsReports
	for rows.Next() {
		r, err := ir.scanDetail(rows)
		if err != nil {
			return nil, err
		}

		resp = append(resp, r)
	}

	return resp, nil
}

func (ir InvoiceReport) scanHead(s pgx.Row) (model.InvoiceReport, error) {
	invoice := model.Invoice{}
	user := model.User{}

	invoiceUpdateAtNull := sql.NullInt64{}
	userUpdateAtNull := sql.NullInt64{}

	err := s.Scan(
		&invoice.ID,
		&invoice.UserID,
		&invoice.PurchaseOrderID,
		&invoice.CreatedAt,
		&invoiceUpdateAtNull,
		&user.ID,
		&user.Email,
		&user.IsAdmin,
		&user.Details,
		&user.CreatedAt,
		&userUpdateAtNull,
	)
	if err != nil {
		return model.InvoiceReport{}, err
	}

	invoice.UpdatedAt = invoiceUpdateAtNull.Int64
	user.UpdatedAt = userUpdateAtNull.Int64

	r := model.InvoiceReport{
		Invoice: invoice,
		User:    user,
	}

	return r, nil
}

func (ir InvoiceReport) scanDetail(s pgx.Row) (model.InvoiceDetailsReport, error) {
	invoiceDetail := model.InvoiceDetail{}
	product := model.Product{}

	invoiceDetailUpdateAtNull := sql.NullInt64{}
	productUpdateAtNull := sql.NullInt64{}

	err := s.Scan(
		&invoiceDetail.ID,
		&invoiceDetail.InvoiceID,
		&invoiceDetail.ProductID,
		&invoiceDetail.Amount,
		&invoiceDetail.UnitPrice,
		&invoiceDetail.CreatedAt,
		&invoiceDetailUpdateAtNull,
		&product.ID,
		&product.ProductName,
		&product.Price,
		&product.Images,
		&product.Description,
		&product.Features,
		&product.CreatedAt,
		&productUpdateAtNull,
	)
	if err != nil {
		return model.InvoiceDetailsReport{}, err
	}

	invoiceDetail.UpdatedAt = invoiceDetailUpdateAtNull.Int64
	product.UpdatedAt = productUpdateAtNull.Int64

	r := model.InvoiceDetailsReport{
		InvoiceDetail: invoiceDetail,
		Product:       product,
	}

	return r, nil
}
