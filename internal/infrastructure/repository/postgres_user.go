package repository

import (
	"database/sql"
	"fmt"

	domain "github.com/k07g/g1/internal/domain/user"
)

const createUsersTableSQL = `
CREATE TABLE IF NOT EXISTS users (
	id   SERIAL PRIMARY KEY,
	name TEXT    NOT NULL,
	age  INTEGER NOT NULL
);`

// PostgreSQLUserRepository はuser.Repositoryインターフェースのメモリ実装
type PostgreSQLUserRepository struct {
	db *sql.DB
}

// NewPostgreSQLUserRepository はDBに接続し、テーブルを初期化してリポジトリを返す
func NewPostgreSQLUserRepository(dsn string) (domain.Repository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("db open: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}
	if _, err := db.Exec(createUsersTableSQL); err != nil {
		return nil, fmt.Errorf("create table: %w", err)
	}
	return &PostgreSQLUserRepository{db: db}, nil
}

func (r *PostgreSQLUserRepository) FindAll() ([]*domain.User, error) {
	rows, err := r.db.Query(`SELECT id, name, age FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		u := &domain.User{}
		if err := rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *PostgreSQLUserRepository) FindByID(id int) (*domain.User, error) {
	u := &domain.User{}
	err := r.db.QueryRow(`SELECT id, name, age FROM users WHERE id = $1`, id).
		Scan(&u.ID, &u.Name, &u.Age)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *PostgreSQLUserRepository) Save(user *domain.User) (*domain.User, error) {
	err := r.db.QueryRow(
		`INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id`,
		user.Name, user.Age,
	).Scan(&user.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *PostgreSQLUserRepository) Update(user *domain.User) (*domain.User, error) {
	result, err := r.db.Exec(
		`UPDATE users SET name=$1, age=$2 WHERE id=$3`,
		user.Name, user.Age, user.ID,
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
	return user, nil
}

func (r *PostgreSQLUserRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM users WHERE id=$1`, id)
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
