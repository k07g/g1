package repository

import (
	"database/sql"
	"fmt"

	domain "github.com/k07g/g1/internal/domain/item"
	_ "github.com/lib/pq"
)

const createTableSQL = `
CREATE TABLE IF NOT EXISTS items (
	id         SERIAL PRIMARY KEY,
	name       TEXT        NOT NULL,
	price      INTEGER     NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL
);`

// PostgreSQLRepository はRepositoryインターフェースのPostgreSQL実装
type PostgreSQLRepository struct {
	db *sql.DB
}

// NewPostgreSQLRepository はDBに接続し、テーブルを初期化してリポジトリを返す
func NewPostgreSQLRepository(dsn string) (domain.Repository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("db open: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}
	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, fmt.Errorf("create table: %w", err)
	}
	return &PostgreSQLRepository{db: db}, nil
}

func (r *PostgreSQLRepository) FindAll() ([]*domain.Item, error) {
	rows, err := r.db.Query(
		`SELECT id, name, price, created_at, updated_at FROM items ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.Item
	for rows.Next() {
		item := &domain.Item{}
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgreSQLRepository) FindByID(id int) (*domain.Item, error) {
	item := &domain.Item{}
	err := r.db.QueryRow(
		`SELECT id, name, price, created_at, updated_at FROM items WHERE id = $1`, id,
	).Scan(&item.ID, &item.Name, &item.Price, &item.CreatedAt, &item.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *PostgreSQLRepository) Save(item *domain.Item) (*domain.Item, error) {
	err := r.db.QueryRow(
		`INSERT INTO items (name, price, created_at, updated_at)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		item.Name, item.Price, item.CreatedAt, item.UpdatedAt,
	).Scan(&item.ID)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *PostgreSQLRepository) Update(item *domain.Item) (*domain.Item, error) {
	result, err := r.db.Exec(
		`UPDATE items SET name=$1, price=$2, updated_at=$3 WHERE id=$4`,
		item.Name, item.Price, item.UpdatedAt, item.ID,
	)
	if err != nil {
		return nil, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, domain.ErrNotFound
	}
	return item, nil
}

func (r *PostgreSQLRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM items WHERE id=$1`, id)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}
