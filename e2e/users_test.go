package e2e

import (
	"fmt"
	"net/http"
	"testing"
)

// ==============================
// GET /api/users
// ==============================

func TestE2E_Users_List(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := do(t, http.MethodGet, srv.URL+"/api/users", nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)

	r := decodeResponse(t, resp)
	assertSuccess(t, r, true)

	users := decodeUsers(t, r)
	// テスト開始時は空
	if len(users) != 0 {
		t.Errorf("users count = %d, want 0", len(users))
	}
}

// ==============================
// POST /api/users
// ==============================

func TestE2E_Users_Create(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	body := map[string]interface{}{"name": "山田太郎", "age": 30}
	resp := do(t, http.MethodPost, srv.URL+"/api/users", body)
	assertStatus(t, resp.StatusCode, http.StatusCreated)

	r := decodeResponse(t, resp)
	assertSuccess(t, r, true)

	u := decodeUser(t, r)
	if u.ID == 0 {
		t.Error("ID should be assigned")
	}
	if u.Name != "山田太郎" {
		t.Errorf("Name = %q, want %q", u.Name, "山田太郎")
	}
	if u.Age != 30 {
		t.Errorf("Age = %d, want 30", u.Age)
	}
}

func TestE2E_Users_Create_ValidationErrors(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	tests := []struct {
		name string
		body map[string]interface{}
	}{
		{"名前が空", map[string]interface{}{"name": "", "age": 20}},
		{"年齢がマイナス", map[string]interface{}{"name": "田中", "age": -1}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp := do(t, http.MethodPost, srv.URL+"/api/users", tt.body)
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
// GET /api/users/{id}
// ==============================

func TestE2E_Users_Get(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	created := decodeUser(t, decodeResponse(t,
		do(t, http.MethodPost, srv.URL+"/api/users",
			map[string]interface{}{"name": "佐藤花子", "age": 25}),
	))

	resp := do(t, http.MethodGet, fmt.Sprintf("%s/api/users/%d", srv.URL, created.ID), nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)

	u := decodeUser(t, decodeResponse(t, resp))
	if u.Name != "佐藤花子" || u.Age != 25 {
		t.Errorf("user = %+v, want Name=佐藤花子 Age=25", u)
	}
}

func TestE2E_Users_Get_NotFound(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := do(t, http.MethodGet, srv.URL+"/api/users/9999", nil)
	assertStatus(t, resp.StatusCode, http.StatusNotFound)
	assertSuccess(t, decodeResponse(t, resp), false)
}

func TestE2E_Users_Get_InvalidID(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := do(t, http.MethodGet, srv.URL+"/api/users/abc", nil)
	assertStatus(t, resp.StatusCode, http.StatusBadRequest)
	assertSuccess(t, decodeResponse(t, resp), false)
}

// ==============================
// PUT /api/users/{id}
// ==============================

func TestE2E_Users_Update(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	created := decodeUser(t, decodeResponse(t,
		do(t, http.MethodPost, srv.URL+"/api/users",
			map[string]interface{}{"name": "更新前", "age": 20}),
	))

	resp := do(t, http.MethodPut,
		fmt.Sprintf("%s/api/users/%d", srv.URL, created.ID),
		map[string]interface{}{"name": "更新後", "age": 21},
	)
	assertStatus(t, resp.StatusCode, http.StatusOK)

	updated := decodeUser(t, decodeResponse(t, resp))
	if updated.Name != "更新後" || updated.Age != 21 {
		t.Errorf("updated = %+v, want Name=更新後 Age=21", updated)
	}
}

func TestE2E_Users_Update_PartialFields(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	created := decodeUser(t, decodeResponse(t,
		do(t, http.MethodPost, srv.URL+"/api/users",
			map[string]interface{}{"name": "元の名前", "age": 30}),
	))

	// 名前のみ更新
	resp := do(t, http.MethodPut,
		fmt.Sprintf("%s/api/users/%d", srv.URL, created.ID),
		map[string]interface{}{"name": "新しい名前"},
	)
	assertStatus(t, resp.StatusCode, http.StatusOK)

	updated := decodeUser(t, decodeResponse(t, resp))
	if updated.Name != "新しい名前" {
		t.Errorf("Name = %q, want %q", updated.Name, "新しい名前")
	}
	// 年齢は変わっていない
	if updated.Age != 30 {
		t.Errorf("Age = %d, want 30 (unchanged)", updated.Age)
	}
}

func TestE2E_Users_Update_NotFound(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := do(t, http.MethodPut, srv.URL+"/api/users/9999",
		map[string]interface{}{"name": "x", "age": 1},
	)
	assertStatus(t, resp.StatusCode, http.StatusNotFound)
	assertSuccess(t, decodeResponse(t, resp), false)
}

// ==============================
// DELETE /api/users/{id}
// ==============================

func TestE2E_Users_Delete(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	created := decodeUser(t, decodeResponse(t,
		do(t, http.MethodPost, srv.URL+"/api/users",
			map[string]interface{}{"name": "削除対象", "age": 18}),
	))

	resp := do(t, http.MethodDelete, fmt.Sprintf("%s/api/users/%d", srv.URL, created.ID), nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	assertSuccess(t, decodeResponse(t, resp), true)

	// 削除後はGETできない
	resp2 := do(t, http.MethodGet, fmt.Sprintf("%s/api/users/%d", srv.URL, created.ID), nil)
	assertStatus(t, resp2.StatusCode, http.StatusNotFound)
}

func TestE2E_Users_Delete_NotFound(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := do(t, http.MethodDelete, srv.URL+"/api/users/9999", nil)
	assertStatus(t, resp.StatusCode, http.StatusNotFound)
	assertSuccess(t, decodeResponse(t, resp), false)
}

// ==============================
// シナリオテスト（CRUD一連の流れ）
// ==============================

func TestE2E_Users_Scenario_CreateAndDelete(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	// 1. 作成
	created := decodeUser(t, decodeResponse(t,
		do(t, http.MethodPost, srv.URL+"/api/users",
			map[string]interface{}{"name": "シナリオユーザー", "age": 40}),
	))

	// 2. 一覧に含まれることを確認
	users := decodeUsers(t, decodeResponse(t,
		do(t, http.MethodGet, srv.URL+"/api/users", nil),
	))
	found := false
	for _, u := range users {
		if u.ID == created.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("created user (ID=%d) not found in list", created.ID)
	}

	// 3. 更新
	updated := decodeUser(t, decodeResponse(t,
		do(t, http.MethodPut,
			fmt.Sprintf("%s/api/users/%d", srv.URL, created.ID),
			map[string]interface{}{"age": 41},
		),
	))
	if updated.Age != 41 {
		t.Errorf("Age = %d, want 41", updated.Age)
	}

	// 4. 削除
	do(t, http.MethodDelete, fmt.Sprintf("%s/api/users/%d", srv.URL, created.ID), nil)

	// 5. 削除後は取得不可
	resp := do(t, http.MethodGet, fmt.Sprintf("%s/api/users/%d", srv.URL, created.ID), nil)
	assertStatus(t, resp.StatusCode, http.StatusNotFound)
}
