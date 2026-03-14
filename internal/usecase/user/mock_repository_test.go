package user_test

import (
	domain "github.com/k07g/g1/internal/domain/user"
)

// mockRepository はテスト用のリポジトリモック
type mockRepository struct {
	users  map[int]*domain.User
	nextID int
	// エラーを注入するフック
	findAllErr  error
	findByIDErr error
	saveErr     error
	updateErr   error
	deleteErr   error
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		users:  make(map[int]*domain.User),
		nextID: 1,
	}
}

func (m *mockRepository) FindAll() ([]*domain.User, error) {
	if m.findAllErr != nil {
		return nil, m.findAllErr
	}
	users := make([]*domain.User, 0, len(m.users))
	for _, u := range m.users {
		copied := *u
		users = append(users, &copied)
	}
	return users, nil
}

func (m *mockRepository) FindByID(id int) (*domain.User, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	u, ok := m.users[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	copied := *u
	return &copied, nil
}

func (m *mockRepository) Save(user *domain.User) (*domain.User, error) {
	if m.saveErr != nil {
		return nil, m.saveErr
	}
	user.ID = m.nextID
	m.nextID++
	copied := *user
	m.users[user.ID] = &copied
	return user, nil
}

func (m *mockRepository) Update(user *domain.User) (*domain.User, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	if _, ok := m.users[user.ID]; !ok {
		return nil, domain.ErrNotFound
	}
	copied := *user
	m.users[user.ID] = &copied
	return user, nil
}

func (m *mockRepository) Delete(id int) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, ok := m.users[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.users, id)
	return nil
}

// seedUser はモックにユーザーを1件追加するヘルパー
func (m *mockRepository) seedUser(name string, age int) *domain.User {
	u := &domain.User{
		ID:   m.nextID,
		Name: name,
		Age:  age,
	}
	m.users[m.nextID] = u
	m.nextID++
	return u
}
