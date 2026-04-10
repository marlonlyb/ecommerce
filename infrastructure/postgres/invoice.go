package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mlbautomation/ProyectoEMLB/model"
)

const (
	iTable = "invoices"
)

var iFields = []string{
	"id",
	"user_id",
	"purchase_order_id",
	"created_at",
	"updated_at",
}

var (
	iPsqlInsert = BuildSQLInsert(iTable, iFields)
)

type Invoice struct {
	db *pgxpool.Pool
}

func NewInvoice(db *pgxpool.Pool) Invoice {
	return Invoice{db: db}
}

func (i Invoice) getTx() (pgx.Tx, error) {
	return i.db.Begin(context.Background())
}

func (i Invoice) Create(m *model.Invoice, ms model.InvoiceDetails) error {
	tx, err := i.getTx()
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		context.Background(),
		iPsqlInsert,
		m.ID,
		m.UserID,
		m.PurchaseOrderID,
		m.CreatedAt,
		Int64ToNull(m.UpdatedAt),
	)
	if err != nil {
		errRollback := tx.Rollback(context.Background())
		if errRollback != nil {
			return fmt.Errorf("%s %w", errRollback, err)
		}

		return err
	}

	err = i.CreateDetailsBulk(tx, ms)
	if err != nil {
		errRollback := tx.Rollback(context.Background())
		if errRollback != nil {
			return fmt.Errorf("%s %w", errRollback, err)
		}

		return err
	}

	//confirmamos la transacción
	errCommit := tx.Commit(context.Background())
	if errCommit != nil {
		return errCommit
	}

	return nil
}
