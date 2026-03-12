package item_test

import (
	"fmt"

	domain "github.com/k07g/g1/internal/domain/item"
)

// mockRepository はテスト用のリポジトリモック
type mockRepository struct {
	items  map[int]*domain.Item
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
		items:  make(map[int]*domain.Item),
		nextID: 1,
	}
}

func (m *mockRepository) FindAll() ([]*domain.Item, error) {
	if m.findAllErr != nil {
		return nil, m.findAllErr
	}
	items := make([]*domain.Item, 0, len(m.items))
	for _, item := range m.items {
		copied := *item
		items = append(items, &copied)
	}
	return items, nil
}

func (m *mockRepository) FindByID(id int) (*domain.Item, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	item, ok := m.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	copied := *item
	return &copied, nil
}

func (m *mockRepository) Save(item *domain.Item) (*domain.Item, error) {
	if m.saveErr != nil {
		return nil, m.saveErr
	}
	item.ID = m.nextID
	m.nextID++
	copied := *item
	m.items[item.ID] = &copied
	return item, nil
}

func (m *mockRepository) Update(item *domain.Item) (*domain.Item, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	if _, ok := m.items[item.ID]; !ok {
		return nil, domain.ErrNotFound
	}
	copied := *item
	m.items[item.ID] = &copied
	return item, nil
}

func (m *mockRepository) Delete(id int) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, ok := m.items[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.items, id)
	return nil
}

// seedItem はモックにアイテムを1件追加するヘルパー
func (m *mockRepository) seedItem(name string, price int) *domain.Item {
	item := &domain.Item{
		ID:    m.nextID,
		Name:  name,
		Price: price,
	}
	m.items[m.nextID] = item
	m.nextID++
	return item
}

// suppress unused import error
var _ = fmt.Sprintf
