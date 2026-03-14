package user

import (
	domain "github.com/k07g/g1/internal/domain/user"
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
	users, err := i.repo.FindAll()
	if err != nil {
		return nil, err
	}
	outputs := make([]*Output, 0, len(users))
	for _, u := range users {
		outputs = append(outputs, toOutput(u))
	}
	return outputs, nil
}

func (i *Interactor) Get(id int) (*Output, error) {
	u, err := i.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return toOutput(u), nil
}

func (i *Interactor) Create(input CreateInput) (*Output, error) {
	u := &domain.User{
		Name: input.Name,
		Age:  input.Age,
	}
	if err := u.Validate(); err != nil {
		return nil, err
	}
	saved, err := i.repo.Save(u)
	if err != nil {
		return nil, err
	}
	return toOutput(saved), nil
}

func (i *Interactor) Update(input UpdateInput) (*Output, error) {
	u, err := i.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		u.Name = input.Name
	}
	if input.Age != 0 {
		u.Age = input.Age
	}
	if err := u.Validate(); err != nil {
		return nil, err
	}
	updated, err := i.repo.Update(u)
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

func toOutput(u *domain.User) *Output {
	return &Output{
		ID:   u.ID,
		Name: u.Name,
		Age:  u.Age,
	}
}
