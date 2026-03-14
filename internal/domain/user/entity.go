package user

import "errors"

// User はドメインエンティティ
type User struct {
	ID   int
	Name string
	Age  int
}

// ドメインエラー定義
var (
	ErrNotFound    = errors.New("user not found")
	ErrInvalidName = errors.New("name must not be empty")
	ErrInvalidAge  = errors.New("age must be 0 or greater")
)

// Validate はエンティティのバリデーションを行う
func (u *User) Validate() error {
	if u.Name == "" {
		return ErrInvalidName
	}
	if u.Age < 0 {
		return ErrInvalidAge
	}
	return nil
}
