package e2e

import (
	"net/http/httptest"

	"github.com/k07g/g1/internal/infrastructure/repository"
	"github.com/k07g/g1/internal/interface/handler"
	usecase "github.com/k07g/g1/internal/usecase/item"
)

// newTestServer は実際の依存関係を組み立てたテスト用サーバーを起動する。
// ユニットテストのモックと異なり、全レイヤーを結合した状態で動作する。
func newTestServer() *httptest.Server {
	repo := repository.NewInMemoryRepository()
	uc := usecase.NewInteractor(repo)
	h := handler.NewItemHandler(uc)

	mux := newServeMux(h)
	return httptest.NewServer(mux)
}
