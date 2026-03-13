package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

// ==============================
// ヘルスチェック
// ==============================

func TestE2E_Health(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp := do(t, http.MethodGet, srv.URL+"/health", nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
}

// ==============================
// GET /api/items
// ==============================

func TestE2E_List(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp := do(t, http.MethodGet, srv.URL+"/api/items", nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)

	r := decodeResponse(t, resp)
	assertSuccess(t, r, true)

	items := decodeItems(t, r)
	// seedデータが3件入っている
	if len(items) != 3 {
		t.Errorf("items count = %d, want 3", len(items))
	}
}

// ==============================
// POST /api/items
// ==============================

func TestE2E_Create(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	body := map[string]interface{}{"name": "新しい本", "price": 2500}
	resp := do(t, http.MethodPost, srv.URL+"/api/items", body)
	assertStatus(t, resp.StatusCode, http.StatusCreated)

	r := decodeResponse(t, resp)
	assertSuccess(t, r, true)

	item := decodeItem(t, r)
	if item.ID == 0 {
		t.Error("ID should be assigned")
	}
	if item.Name != "新しい本" {
		t.Errorf("Name = %q, want %q", item.Name, "新しい本")
	}
	if item.Price != 2500 {
		t.Errorf("Price = %d, want 2500", item.Price)
	}
	if item.CreatedAt == "" || item.UpdatedAt == "" {
		t.Error("CreatedAt/UpdatedAt should be set")
	}
}

func TestE2E_Create_ValidationErrors(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	tests := []struct {
		name string
		body map[string]interface{}
	}{
		{"名前が空", map[string]interface{}{"name": "", "price": 100}},
		{"価格がマイナス", map[string]interface{}{"name": "商品", "price": -1}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp := do(t, http.MethodPost, srv.URL+"/api/items", tt.body)
			assertStatus(t, resp.StatusCode, http.StatusBadRequest)
			r := decodeResponse(t, resp)
			assertSuccess(t, r, false)
			if r.Error == "" {
				t.Error("error message should not be empty")
			}
		})
	}
}

// ==============================
// GET /api/items/{id}
// ==============================

func TestE2E_Get(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	// まず作成
	created := decodeItem(t, decodeResponse(t,
		do(t, http.MethodPost, srv.URL+"/api/items",
			map[string]interface{}{"name": "取得テスト", "price": 1000}),
	))

	resp := do(t, http.MethodGet, fmt.Sprintf("%s/api/items/%d", srv.URL, created.ID), nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)

	item := decodeItem(t, decodeResponse(t, resp))
	if item.Name != "取得テスト" || item.Price != 1000 {
		t.Errorf("item = %+v, want Name=取得テスト Price=1000", item)
	}
}

func TestE2E_Get_NotFound(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp := do(t, http.MethodGet, srv.URL+"/api/items/9999", nil)
	assertStatus(t, resp.StatusCode, http.StatusNotFound)
	assertSuccess(t, decodeResponse(t, resp), false)
}

func TestE2E_Get_InvalidID(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp := do(t, http.MethodGet, srv.URL+"/api/items/abc", nil)
	assertStatus(t, resp.StatusCode, http.StatusBadRequest)
	assertSuccess(t, decodeResponse(t, resp), false)
}

// ==============================
// PUT /api/items/{id}
// ==============================

func TestE2E_Update(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	// 作成
	created := decodeItem(t, decodeResponse(t,
		do(t, http.MethodPost, srv.URL+"/api/items",
			map[string]interface{}{"name": "更新前", "price": 100}),
	))

	// UpdatedAt の差分を確認するために1秒待機
	time.Sleep(1 * time.Second)

	// 更新
	resp := do(t, http.MethodPut,
		fmt.Sprintf("%s/api/items/%d", srv.URL, created.ID),
		map[string]interface{}{"name": "更新後", "price": 999},
	)
	assertStatus(t, resp.StatusCode, http.StatusOK)

	updated := decodeItem(t, decodeResponse(t, resp))
	if updated.Name != "更新後" || updated.Price != 999 {
		t.Errorf("updated = %+v, want Name=更新後 Price=999", updated)
	}
	if updated.UpdatedAt == created.UpdatedAt {
		t.Error("UpdatedAt should be changed after update")
	}
}

func TestE2E_Update_PartialFields(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	created := decodeItem(t, decodeResponse(t,
		do(t, http.MethodPost, srv.URL+"/api/items",
			map[string]interface{}{"name": "元の名前", "price": 500}),
	))

	// 名前のみ更新
	resp := do(t, http.MethodPut,
		fmt.Sprintf("%s/api/items/%d", srv.URL, created.ID),
		map[string]interface{}{"name": "新しい名前"},
	)
	assertStatus(t, resp.StatusCode, http.StatusOK)

	updated := decodeItem(t, decodeResponse(t, resp))
	if updated.Name != "新しい名前" {
		t.Errorf("Name = %q, want %q", updated.Name, "新しい名前")
	}
	// 価格は変わっていない
	if updated.Price != 500 {
		t.Errorf("Price = %d, want 500 (unchanged)", updated.Price)
	}
}

func TestE2E_Update_NotFound(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp := do(t, http.MethodPut, srv.URL+"/api/items/9999",
		map[string]interface{}{"name": "x", "price": 1},
	)
	assertStatus(t, resp.StatusCode, http.StatusNotFound)
	assertSuccess(t, decodeResponse(t, resp), false)
}

// ==============================
// DELETE /api/items/{id}
// ==============================

func TestE2E_Delete(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	// 作成
	created := decodeItem(t, decodeResponse(t,
		do(t, http.MethodPost, srv.URL+"/api/items",
			map[string]interface{}{"name": "削除対象", "price": 100}),
	))

	// 削除
	resp := do(t, http.MethodDelete, fmt.Sprintf("%s/api/items/%d", srv.URL, created.ID), nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	assertSuccess(t, decodeResponse(t, resp), true)

	// 削除後はGETできない
	resp2 := do(t, http.MethodGet, fmt.Sprintf("%s/api/items/%d", srv.URL, created.ID), nil)
	assertStatus(t, resp2.StatusCode, http.StatusNotFound)
}

func TestE2E_Delete_NotFound(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp := do(t, http.MethodDelete, srv.URL+"/api/items/9999", nil)
	assertStatus(t, resp.StatusCode, http.StatusNotFound)
	assertSuccess(t, decodeResponse(t, resp), false)
}

// ==============================
// シナリオテスト（CRUD一連の流れ）
// ==============================

func TestE2E_Scenario_CreateAndDelete(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	// 1. 作成
	created := decodeItem(t, decodeResponse(t,
		do(t, http.MethodPost, srv.URL+"/api/items",
			map[string]interface{}{"name": "シナリオ商品", "price": 1500}),
	))

	// 2. 一覧に含まれることを確認
	items := decodeItems(t, decodeResponse(t,
		do(t, http.MethodGet, srv.URL+"/api/items", nil),
	))
	found := false
	for _, item := range items {
		if item.ID == created.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("created item (ID=%d) not found in list", created.ID)
	}

	// 3. 更新
	updated := decodeItem(t, decodeResponse(t,
		do(t, http.MethodPut,
			fmt.Sprintf("%s/api/items/%d", srv.URL, created.ID),
			map[string]interface{}{"price": 2000},
		),
	))
	if updated.Price != 2000 {
		t.Errorf("Price = %d, want 2000", updated.Price)
	}

	// 4. 削除
	do(t, http.MethodDelete, fmt.Sprintf("%s/api/items/%d", srv.URL, created.ID), nil)

	// 5. 削除後は取得不可
	resp := do(t, http.MethodGet, fmt.Sprintf("%s/api/items/%d", srv.URL, created.ID), nil)
	assertStatus(t, resp.StatusCode, http.StatusNotFound)
}

// ==============================
// メソッド不正
// ==============================

func TestE2E_MethodNotAllowed(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp := do(t, http.MethodPatch, srv.URL+"/api/items", map[string]interface{}{"name": "x"})
	assertStatus(t, resp.StatusCode, http.StatusMethodNotAllowed)
}
