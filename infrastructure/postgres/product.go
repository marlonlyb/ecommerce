package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mlbautomation/ProyectoEMLB/model"
)

const pTable = "products"

var pFields = []string{
	"id",
	"product_name",
	"price",
	"images",
	"description",
	"features",
	"created_at",
	"updated_at",
}

var (
	pPsqlInsert = BuildSQLInsert(pTable, pFields)
	pPsqlUpdate = BuildSQLUpdatedByID(pTable, pFields)
	pPsqlDelete = BuildSQLDelete(pTable)
	pPsqlGetAll = BuildSQLSelect(pTable, pFields)
)

type Product struct {
	db *pgxpool.Pool
}

func NewProduct(db *pgxpool.Pool) Product {
	return Product{db: db}
}

func (p Product) Create(m *model.Product) error {
	_, err := p.db.Exec(
		context.Background(),
		pPsqlInsert,
		m.ID,
		m.ProductName,
		m.Price,
		m.Images,
		m.Description,
		m.Features,
		m.CreatedAt,
		Int64ToNull(m.UpdatedAt),
	)

	if err != nil {
		return err
	}
	return nil
}

func (p Product) Update(m *model.Product) error {
	_, err := p.db.Exec(
		context.Background(),
		pPsqlUpdate,
		m.ProductName,
		m.Price,
		m.Images,
		m.Description,
		m.Features,
		m.UpdatedAt,
		m.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p Product) Delete(ID uuid.UUID) error {
	_, err := p.db.Exec(
		context.Background(),
		pPsqlDelete,
		ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p Product) GetByID(ID uuid.UUID) (model.Product, error) {
	query := pPsqlGetAll + " WHERE id = $1"
	row := p.db.QueryRow(
		context.Background(),
		query,
		ID,
	)
	return p.scanRow(row)
}

func (p Product) GetAll() (model.Products, error) {
	rows, err := p.db.Query(
		context.Background(),
		pPsqlGetAll,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var ms model.Products

	for rows.Next() {
		m, err := p.scanRow(rows)
		if err != nil {
			return nil, err
		}

		ms = append(ms, m)
	}

	return ms, nil
}

func (p Product) scanRow(s pgx.Row) (model.Product, error) {

	var m model.Product

	updateAtNull := sql.NullInt64{}

	err := s.Scan(
		&m.ID,
		&m.ProductName,
		&m.Price,
		&m.Images,
		&m.Description,
		&m.Features,
		&m.CreatedAt,
		&updateAtNull,
	)
	if err != nil {
		return m, err
	}

	m.UpdatedAt = updateAtNull.Int64

	return m, nil
}
