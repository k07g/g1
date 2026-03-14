package user

// --- Input DTO（ユースケースへの入力）---

type CreateInput struct {
	Name string
	Age  int
}

type UpdateInput struct {
	ID   int
	Name string
	Age  int
}

// --- Output DTO（ユースケースからの出力）---

type Output struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// --- ユースケースインターフェース ---

type UseCase interface {
	List() ([]*Output, error)
	Get(id int) (*Output, error)
	Create(input CreateInput) (*Output, error)
	Update(input UpdateInput) (*Output, error)
	Delete(id int) error
}
