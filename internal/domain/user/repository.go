package user

// Repository はデータ永続化の抽象インターフェース（依存性逆転の原則）
type Repository interface {
	FindAll() ([]*User, error)
	FindByID(id int) (*User, error)
	Save(user *User) (*User, error)
	Update(user *User) (*User, error)
	Delete(id int) error
}
