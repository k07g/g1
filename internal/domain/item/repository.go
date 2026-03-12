package item

// Repository はデータ永続化の抽象インターフェース（依存性逆転の原則）
type Repository interface {
	FindAll() ([]*Item, error)
	FindByID(id int) (*Item, error)
	Save(item *Item) (*Item, error)
	Update(item *Item) (*Item, error)
	Delete(id int) error
}
