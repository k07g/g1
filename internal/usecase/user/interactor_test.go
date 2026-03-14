package user_test

import (
	"errors"
	"testing"

	domain "github.com/k07g/g1/internal/domain/user"
	usecase "github.com/k07g/g1/internal/usecase/user"
)

func newInteractor(repo *mockRepository) usecase.UseCase {
	return usecase.NewInteractor(repo)
}

// ---- List ----

func TestInteractor_List(t *testing.T) {
	t.Parallel()
	repo := newMockRepository()
	repo.seedUser("ユーザーA", 20)
	repo.seedUser("ユーザーB", 30)

	uc := newInteractor(repo)
	outputs, err := uc.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(outputs) != 2 {
		t.Errorf("List() len = %d, want 2", len(outputs))
	}
}

func TestInteractor_List_RepoError(t *testing.T) {
	t.Parallel()
	repo := newMockRepository()
	repo.findAllErr = errors.New("db error")

	uc := newInteractor(repo)
	_, err := uc.List()
	if err == nil {
		t.Error("List() should return error when repo fails")
	}
}

// ---- Get ----

func TestInteractor_Get(t *testing.T) {
	t.Parallel()
	repo := newMockRepository()
	seeded := repo.seedUser("山田太郎", 25)

	uc := newInteractor(repo)
	output, err := uc.Get(seeded.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if output.Name != "山田太郎" || output.Age != 25 {
		t.Errorf("Get() = %+v, want Name=山田太郎 Age=25", output)
	}
}

func TestInteractor_Get_NotFound(t *testing.T) {
	t.Parallel()
	uc := newInteractor(newMockRepository())

	_, err := uc.Get(9999)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Get() error = %v, want ErrNotFound", err)
	}
}

// ---- Create ----

func TestInteractor_Create(t *testing.T) {
	t.Parallel()
	uc := newInteractor(newMockRepository())

	output, err := uc.Create(usecase.CreateInput{Name: "新ユーザー", Age: 18})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if output.ID == 0 {
		t.Error("Create() ID should be assigned")
	}
	if output.Name != "新ユーザー" || output.Age != 18 {
		t.Errorf("Create() = %+v, want Name=新ユーザー Age=18", output)
	}
}

func TestInteractor_Create_ValidationError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   usecase.CreateInput
		wantErr error
	}{
		{
			name:    "名前が空",
			input:   usecase.CreateInput{Name: "", Age: 20},
			wantErr: domain.ErrInvalidName,
		},
		{
			name:    "年齢がマイナス",
			input:   usecase.CreateInput{Name: "ユーザー", Age: -1},
			wantErr: domain.ErrInvalidAge,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			uc := newInteractor(newMockRepository())
			_, err := uc.Create(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Create() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

// ---- Update ----

func TestInteractor_Update(t *testing.T) {
	t.Parallel()
	repo := newMockRepository()
	seeded := repo.seedUser("旧名前", 20)

	uc := newInteractor(repo)
	output, err := uc.Update(usecase.UpdateInput{ID: seeded.ID, Name: "新名前", Age: 21})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if output.Name != "新名前" || output.Age != 21 {
		t.Errorf("Update() = %+v, want Name=新名前 Age=21", output)
	}
}

func TestInteractor_Update_NotFound(t *testing.T) {
	t.Parallel()
	uc := newInteractor(newMockRepository())

	_, err := uc.Update(usecase.UpdateInput{ID: 9999, Name: "x"})
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Update() error = %v, want ErrNotFound", err)
	}
}

func TestInteractor_Update_PartialFields(t *testing.T) {
	t.Parallel()
	repo := newMockRepository()
	seeded := repo.seedUser("元の名前", 30)

	uc := newInteractor(repo)
	// Nameのみ更新
	output, err := uc.Update(usecase.UpdateInput{ID: seeded.ID, Name: "新しい名前"})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	// Ageは元の値のまま
	if output.Age != 30 {
		t.Errorf("Update() Age = %d, want 30 (unchanged)", output.Age)
	}
}

// ---- Delete ----

func TestInteractor_Delete(t *testing.T) {
	t.Parallel()
	repo := newMockRepository()
	seeded := repo.seedUser("削除対象", 25)

	uc := newInteractor(repo)
	if err := uc.Delete(seeded.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	// 削除後はGetできない
	_, err := uc.Get(seeded.ID)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("after Delete, Get() error = %v, want ErrNotFound", err)
	}
}

func TestInteractor_Delete_NotFound(t *testing.T) {
	t.Parallel()
	uc := newInteractor(newMockRepository())

	err := uc.Delete(9999)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Delete() error = %v, want ErrNotFound", err)
	}
}
