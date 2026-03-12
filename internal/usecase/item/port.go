package item

// --- Input DTO（ユースケースへの入力）---

type CreateInput struct {
	Name  string
	Price int
}

type UpdateInput struct {
	ID    int
	Name  string
	Price int
}

// --- Output DTO（ユースケースからの出力）---

type Output struct {
	ID        int
	Name      string
	Price     int
	CreatedAt string
	UpdatedAt string
}

// --- ユースケースインターフェース ---

type UseCase interface {
	List() ([]*Output, error)
	Get(id int) (*Output, error)
	Create(input CreateInput) (*Output, error)
	Update(input UpdateInput) (*Output, error)
	Delete(id int) error
}
