package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mlbautomation/ProyectoEMLB/model"
)

const poTable = "purchase_orders"

var poFields = []string{
	"id",
	"user_id",
	"products",
	"created_at",
	"updated_at",
}

var (
	poPsqlInsert = BuildSQLInsert(poTable, poFields)
	poPsqlGetAll = BuildSQLSelect(poTable, poFields)
)

type PurchaseOrder struct {
	db *pgxpool.Pool
}

// New returns a new PurchaseOrder storage
func NewPurchaseOrder(db *pgxpool.Pool) *PurchaseOrder {
	return &PurchaseOrder{db: db}
}

// Create creates a model.PurchaseOrder
func (p *PurchaseOrder) Create(m *model.PurchaseOrder) error {
	_, err := p.db.Exec(
		context.Background(),
		poPsqlInsert,
		m.ID,
		m.UserID,
		m.Products,
		m.CreatedAt,
		Int64ToNull(m.UpdatedAt),
	)

	if err != nil {
		return err
	}

	return nil
}

func (p *PurchaseOrder) GetByID(ID uuid.UUID) (model.PurchaseOrder, error) {
	query := poPsqlGetAll + " WHERE id = $1"
	row := p.db.QueryRow(
		context.Background(),
		query,
		ID,
	)

	return p.scanRow(row)
}

func (p *PurchaseOrder) scanRow(s pgx.Row) (model.PurchaseOrder, error) {
	m := model.PurchaseOrder{}

	updatedAtNull := sql.NullInt64{}

	err := s.Scan(
		&m.ID,
		&m.UserID,
		&m.Products,
		&m.CreatedAt,
		&updatedAtNull,
	)
	if err != nil {
		return m, err
	}

	m.UpdatedAt = updatedAtNull.Int64

	return m, nil
}
