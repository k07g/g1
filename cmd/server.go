package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/k07g/g1/internal/infrastructure/repository"
	"github.com/k07g/g1/internal/interface/handler"
	usecase "github.com/k07g/g1/internal/usecase/item"
)

func main() {
	// ---- 依存性の注入（外側から内側へ） ----
	//
	//  [infrastructure]  →  [usecase]  →  [interface/handler]
	//   Repo                 Interactor      ItemHandler
	//       ↓                    ↓
	//  domain.Repository   usecase.UseCase  （インターフェース経由）

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL が設定されていません")
	}
	repo, err := repository.NewPostgreSQLRepository(dsn)
	if err != nil {
		log.Fatalf("PostgreSQL接続エラー: %v", err)
	}
	log.Println("🐘 PostgreSQLリポジトリを使用")
	uc := usecase.NewInteractor(repo)
	itemHandler := handler.NewItemHandler(uc)

	// ---- ルーティング ----
	mux := http.NewServeMux()
	mux.Handle("/api/items", itemHandler)
	mux.Handle("/api/items/", itemHandler)
	mux.HandleFunc("/health", healthHandler)

	// ---- サーバー起動 ----
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("🚀 サーバー起動: http://localhost:8080")
	log.Println("📋 エンドポイント一覧:")
	log.Println("   GET    /api/items       - 一覧取得")
	log.Println("   POST   /api/items       - 新規作成")
	log.Println("   GET    /api/items/{id}  - 1件取得")
	log.Println("   PUT    /api/items/{id}  - 更新")
	log.Println("   DELETE /api/items/{id}  - 削除")
	log.Println("   GET    /health          - ヘルスチェック")

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("サーバーエラー: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
