package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

// response はAPIレスポンスの共通構造体
type response struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// itemResponse は単一アイテムのレスポンス
type itemResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// userResponse は単一ユーザーのレスポンス
type userResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// decodeUser は response.Data を userResponse にデコードする
func decodeUser(t *testing.T, r response) userResponse {
	t.Helper()
	var u userResponse
	if err := json.Unmarshal(r.Data, &u); err != nil {
		t.Fatalf("decode user: %v", err)
	}
	return u
}

// decodeUsers は response.Data を []userResponse にデコードする
func decodeUsers(t *testing.T, r response) []userResponse {
	t.Helper()
	var users []userResponse
	if err := json.Unmarshal(r.Data, &users); err != nil {
		t.Fatalf("decode users: %v", err)
	}
	return users
}

// do はHTTPリクエストを送信してレスポンスを返す
func do(t *testing.T, method, url string, body interface{}) *http.Response {
	t.Helper()
	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request: %v", err)
		}
		r = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, url, r)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	return resp
}

// decodeResponse はレスポンスボディを response 構造体にデコードする
func decodeResponse(t *testing.T, resp *http.Response) response {
	t.Helper()
	defer resp.Body.Close()
	var r response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return r
}

// decodeItem は response.Data を itemResponse にデコードする
func decodeItem(t *testing.T, r response) itemResponse {
	t.Helper()
	var item itemResponse
	if err := json.Unmarshal(r.Data, &item); err != nil {
		t.Fatalf("decode item: %v", err)
	}
	return item
}

// decodeItems は response.Data を []itemResponse にデコードする
func decodeItems(t *testing.T, r response) []itemResponse {
	t.Helper()
	var items []itemResponse
	if err := json.Unmarshal(r.Data, &items); err != nil {
		t.Fatalf("decode items: %v", err)
	}
	return items
}

// assertStatus はステータスコードを検証する
func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("status = %d, want %d", got, want)
	}
}

// assertSuccess はレスポンスの success フィールドを検証する
func assertSuccess(t *testing.T, r response, want bool) {
	t.Helper()
	if r.Success != want {
		t.Errorf("success = %v, want %v (error: %s)", r.Success, want, r.Error)
	}
}
