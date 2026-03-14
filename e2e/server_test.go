package e2e

import (
	"database/sql"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/k07g/g1/internal/infrastructure/repository"
	"github.com/k07g/g1/internal/interface/handler"
	itemUsecase "github.com/k07g/g1/internal/usecase/item"
	userUsecase "github.com/k07g/g1/internal/usecase/user"
	_ "github.com/lib/pq"
)

// newTestServer は実際の依存関係を組み立てたテスト用サーバーを起動する。
// DATABASE_URL 環境変数が必須。テスト前にテーブルをリセット＆シードする。
func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Fatal("DATABASE_URL が設定されていません")
	}

	// CREATE TABLE IF NOT EXISTS を実行するので先にリポジトリを初期化する
	itemRepo, err := repository.NewPostgreSQLRepository(dsn)
	if err != nil {
		t.Fatalf("connect to postgres(items): %v", err)
	}
	userRepo, err := repository.NewPostgreSQLUserRepository(dsn)
	if err != nil {
		t.Fatalf("connect to postgres(users): %v", err)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if _, err := db.Exec(`TRUNCATE TABLE items, users RESTART IDENTITY`); err != nil {
		db.Close()
		t.Fatalf("truncate tables: %v", err)
	}
	if err := seedItems(t, db); err != nil {
		db.Close()
		t.Fatalf("seed items: %v", err)
	}
	db.Close()

	itemHandler := handler.NewItemHandler(itemUsecase.NewInteractor(itemRepo))
	userHandler := handler.NewUserHandler(userUsecase.NewInteractor(userRepo))

	mux := newServeMux(itemHandler, userHandler)
	return httptest.NewServer(mux)
}

var itemSeeds = []struct {
	name  string
	price int
}{
	{"Go言語入門", 3000},
	{"Clean Architecture", 4500},
	{"Docker実践ガイド", 3800},
}

func seedItems(t *testing.T, db *sql.DB) error {
	t.Helper()
	now := time.Now()
	for _, s := range itemSeeds {
		if _, err := db.Exec(
			`INSERT INTO items (name, price, created_at, updated_at) VALUES ($1, $2, $3, $4)`,
			s.name, s.price, now, now,
		); err != nil {
			return err
		}
	}
	return nil
}
