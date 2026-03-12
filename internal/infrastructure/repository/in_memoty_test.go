package repository_test

import (
	"testing"

	domain "github.com/k07g/g1/internal/domain/item"
	"github.com/k07g/g1/internal/infrastructure/repository"
)

func newRepo(t *testing.T) domain.Repository {
	t.Helper()
	return repository.NewInMemoryRepository()
}

func TestInMemoryRepository_FindAll(t *testing.T) {
	t.Parallel()
	repo := newRepo(t)

	items, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll() error = %v", err)
	}
	// seed データが3件入っている
	if len(items) != 3 {
		t.Errorf("FindAll() len = %d, want 3", len(items))
	}
}

func TestInMemoryRepository_Save_And_FindByID(t *testing.T) {
	t.Parallel()
	repo := newRepo(t)

	item := &domain.Item{Name: "新商品", Price: 500}
	saved, err := repo.Save(item)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	if saved.ID == 0 {
		t.Error("Save() ID should be assigned")
	}

	found, err := repo.FindByID(saved.ID)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}
	if found.Name != "新商品" || found.Price != 500 {
		t.Errorf("FindByID() = %+v, want Name=新商品 Price=500", found)
	}
}

func TestInMemoryRepository_FindByID_NotFound(t *testing.T) {
	t.Parallel()
	repo := newRepo(t)

	_, err := repo.FindByID(9999)
	if err != domain.ErrNotFound {
		t.Errorf("FindByID() error = %v, want ErrNotFound", err)
	}
}

func TestInMemoryRepository_Update(t *testing.T) {
	t.Parallel()
	repo := newRepo(t)

	item := &domain.Item{Name: "更新前", Price: 100}
	saved, _ := repo.Save(item)

	saved.Name = "更新後"
	saved.Price = 200
	updated, err := repo.Update(saved)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if updated.Name != "更新後" || updated.Price != 200 {
		t.Errorf("Update() = %+v, want Name=更新後 Price=200", updated)
	}
}

func TestInMemoryRepository_Update_NotFound(t *testing.T) {
	t.Parallel()
	repo := newRepo(t)

	_, err := repo.Update(&domain.Item{ID: 9999, Name: "x", Price: 0})
	if err != domain.ErrNotFound {
		t.Errorf("Update() error = %v, want ErrNotFound", err)
	}
}

func TestInMemoryRepository_Delete(t *testing.T) {
	t.Parallel()
	repo := newRepo(t)

	item := &domain.Item{Name: "削除対象", Price: 100}
	saved, _ := repo.Save(item)

	if err := repo.Delete(saved.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	_, err := repo.FindByID(saved.ID)
	if err != domain.ErrNotFound {
		t.Errorf("after Delete, FindByID() error = %v, want ErrNotFound", err)
	}
}

func TestInMemoryRepository_Delete_NotFound(t *testing.T) {
	t.Parallel()
	repo := newRepo(t)

	err := repo.Delete(9999)
	if err != domain.ErrNotFound {
		t.Errorf("Delete() error = %v, want ErrNotFound", err)
	}
}
