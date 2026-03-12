package item

import (
	"time"

	domain "github.com/k07g/g1/internal/domain/item"
)

// Interactor はユースケースの実装（ビジネスロジックの中心）
type Interactor struct {
	repo domain.Repository
}

// NewInteractor はInteractorを生成する（DIコンストラクタ）
func NewInteractor(repo domain.Repository) UseCase {
	return &Interactor{repo: repo}
}

func (i *Interactor) List() ([]*Output, error) {
	items, err := i.repo.FindAll()
	if err != nil {
		return nil, err
	}
	outputs := make([]*Output, 0, len(items))
	for _, item := range items {
		outputs = append(outputs, toOutput(item))
	}
	return outputs, nil
}

func (i *Interactor) Get(id int) (*Output, error) {
	item, err := i.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return toOutput(item), nil
}

func (i *Interactor) Create(input CreateInput) (*Output, error) {
	now := time.Now()
	item := &domain.Item{
		Name:      input.Name,
		Price:     input.Price,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := item.Validate(); err != nil {
		return nil, err
	}
	saved, err := i.repo.Save(item)
	if err != nil {
		return nil, err
	}
	return toOutput(saved), nil
}

func (i *Interactor) Update(input UpdateInput) (*Output, error) {
	item, err := i.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		item.Name = input.Name
	}
	if input.Price != 0 {
		item.Price = input.Price
	}
	item.UpdatedAt = time.Now()

	if err := item.Validate(); err != nil {
		return nil, err
	}
	updated, err := i.repo.Update(item)
	if err != nil {
		return nil, err
	}
	return toOutput(updated), nil
}

func (i *Interactor) Delete(id int) error {
	if _, err := i.repo.FindByID(id); err != nil {
		return err
	}
	return i.repo.Delete(id)
}

// toOutput はドメインエンティティをOutputDTOに変換する
func toOutput(item *domain.Item) *Output {
	return &Output{
		ID:        item.ID,
		Name:      item.Name,
		Price:     item.Price,
		CreatedAt: item.CreatedAt.Format(time.RFC3339),
		UpdatedAt: item.UpdatedAt.Format(time.RFC3339),
	}
}
