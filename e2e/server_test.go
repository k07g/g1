package e2e

import (
	"database/sql"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/k07g/g1/internal/infrastructure/repository"
	"github.com/k07g/g1/internal/interface/handler"
	usecase "github.com/k07g/g1/internal/usecase/item"
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

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if _, err := db.Exec(`TRUNCATE TABLE items RESTART IDENTITY`); err != nil {
		db.Close()
		t.Fatalf("truncate items: %v", err)
	}
	if err := seedPostgres(t, db); err != nil {
		db.Close()
		t.Fatalf("seed items: %v", err)
	}
	db.Close()

	repo, err := repository.NewPostgreSQLRepository(dsn)
	if err != nil {
		t.Fatalf("connect to postgres: %v", err)
	}

	uc := usecase.NewInteractor(repo)
	h := handler.NewItemHandler(uc)

	mux := newServeMux(h)
	return httptest.NewServer(mux)
}

var seedItems = []struct {
	name  string
	price int
}{
	{"Go言語入門", 3000},
	{"Clean Architecture", 4500},
	{"Docker実践ガイド", 3800},
}

func seedPostgres(t *testing.T, db *sql.DB) error {
	t.Helper()
	now := time.Now()
	for _, s := range seedItems {
		if _, err := db.Exec(
			`INSERT INTO items (name, price, created_at, updated_at) VALUES ($1, $2, $3, $4)`,
			s.name, s.price, now, now,
		); err != nil {
			return err
		}
	}
	return nil
}
