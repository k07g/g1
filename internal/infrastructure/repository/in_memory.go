package repository

import (
	"sync"
	"time"

	domain "github.com/k07g/g1/internal/domain/item"
)

// InMemoryRepository はRepositoryインターフェースのインメモリ実装
// PostgreSQLなど他のDBに切り替える場合はこのファイルだけ差し替えればよい
type InMemoryRepository struct {
	mu     sync.RWMutex
	items  map[int]*domain.Item
	nextID int
}

func NewInMemoryRepository() domain.Repository {
	r := &InMemoryRepository{
		items:  make(map[int]*domain.Item),
		nextID: 1,
	}
	r.seed()
	return r
}

// seed はサンプルデータを投入する
func (r *InMemoryRepository) seed() {
	now := time.Now()
	seeds := []domain.Item{
		{Name: "Go言語入門", Price: 3000, CreatedAt: now, UpdatedAt: now},
		{Name: "Clean Architecture", Price: 4500, CreatedAt: now, UpdatedAt: now},
		{Name: "Docker実践ガイド", Price: 3800, CreatedAt: now, UpdatedAt: now},
	}
	for _, s := range seeds {
		item := s
		item.ID = r.nextID
		r.items[r.nextID] = &item
		r.nextID++
	}
}

func (r *InMemoryRepository) FindAll() ([]*domain.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]*domain.Item, 0, len(r.items))
	for _, item := range r.items {
		copied := *item
		items = append(items, &copied)
	}
	return items, nil
}

func (r *InMemoryRepository) FindByID(id int) (*domain.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	copied := *item
	return &copied, nil
}

func (r *InMemoryRepository) Save(item *domain.Item) (*domain.Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item.ID = r.nextID
	r.nextID++
	copied := *item
	r.items[item.ID] = &copied
	return item, nil
}

func (r *InMemoryRepository) Update(item *domain.Item) (*domain.Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[item.ID]; !ok {
		return nil, domain.ErrNotFound
	}
	copied := *item
	r.items[item.ID] = &copied
	return item, nil
}

func (r *InMemoryRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[id]; !ok {
		return domain.ErrNotFound
	}
	delete(r.items, id)
	return nil
}
