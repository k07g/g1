package item

import (
	"errors"
	"time"
)

// Item はドメインエンティティ
type Item struct {
	ID        int
	Name      string
	Price     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ドメインエラー定義
var (
	ErrNotFound     = errors.New("item not found")
	ErrInvalidName  = errors.New("name must not be empty")
	ErrInvalidPrice = errors.New("price must be 0 or greater")
)

// Validate はエンティティのバリデーションを行う
func (i *Item) Validate() error {
	if i.Name == "" {
		return ErrInvalidName
	}
	if i.Price < 0 {
		return ErrInvalidPrice
	}
	return nil
}
