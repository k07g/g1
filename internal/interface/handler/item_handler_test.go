package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	domain "github.com/k07g/g1/internal/domain/item"
	"github.com/k07g/g1/internal/interface/handler"
	usecase "github.com/k07g/g1/internal/usecase/item"
)

// ---- モックユースケース ----

type mockUseCase struct {
	listOut   []*usecase.Output
	getOut    *usecase.Output
	createOut *usecase.Output
	updateOut *usecase.Output
	err       error
}

func (m *mockUseCase) List() ([]*usecase.Output, error)   { return m.listOut, m.err }
func (m *mockUseCase) Get(_ int) (*usecase.Output, error) { return m.getOut, m.err }
func (m *mockUseCase) Create(_ usecase.CreateInput) (*usecase.Output, error) {
	return m.createOut, m.err
}
func (m *mockUseCase) Update(_ usecase.UpdateInput) (*usecase.Output, error) {
	return m.updateOut, m.err
}
func (m *mockUseCase) Delete(_ int) error { return m.err }

// ---- ヘルパー ----

type apiResp struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   string          `json:"error"`
}

func doRequest(t *testing.T, h http.Handler, method, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}

func decodeResp(t *testing.T, rr *httptest.ResponseRecorder) apiResp {
	t.Helper()
	var resp apiResp
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return resp
}

// ---- GET /api/items ----

func TestItemHandler_List_OK(t *testing.T) {
	t.Parallel()
	uc := &mockUseCase{
		listOut: []*usecase.Output{
			{ID: 1, Name: "商品A", Price: 100},
			{ID: 2, Name: "商品B", Price: 200},
		},
	}
	h := handler.NewItemHandler(uc)

	rr := doRequest(t, h, http.MethodGet, "/api/items", nil)
	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	resp := decodeResp(t, rr)
	if !resp.Success {
		t.Errorf("success = false, want true")
	}
}

func TestItemHandler_List_InternalError(t *testing.T) {
	t.Parallel()
	uc := &mockUseCase{err: errors.New("db error")}
	h := handler.NewItemHandler(uc)

	rr := doRequest(t, h, http.MethodGet, "/api/items", nil)
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusInternalServerError)
	}
}

// ---- GET /api/items/{id} ----

func TestItemHandler_Get_OK(t *testing.T) {
	t.Parallel()
	uc := &mockUseCase{getOut: &usecase.Output{ID: 1, Name: "商品A", Price: 100}}
	h := handler.NewItemHandler(uc)

	rr := doRequest(t, h, http.MethodGet, "/api/items/1", nil)
	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestItemHandler_Get_NotFound(t *testing.T) {
	t.Parallel()
	uc := &mockUseCase{err: domain.ErrNotFound}
	h := handler.NewItemHandler(uc)

	rr := doRequest(t, h, http.MethodGet, "/api/items/9999", nil)
	if rr.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}

func TestItemHandler_Get_InvalidID(t *testing.T) {
	t.Parallel()
	h := handler.NewItemHandler(&mockUseCase{})

	rr := doRequest(t, h, http.MethodGet, "/api/items/abc", nil)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

// ---- POST /api/items ----

func TestItemHandler_Create_OK(t *testing.T) {
	t.Parallel()
	uc := &mockUseCase{createOut: &usecase.Output{ID: 1, Name: "新商品", Price: 500}}
	h := handler.NewItemHandler(uc)

	rr := doRequest(t, h, http.MethodPost, "/api/items", map[string]interface{}{
		"name": "新商品", "price": 500,
	})
	if rr.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusCreated)
	}
	resp := decodeResp(t, rr)
	if !resp.Success {
		t.Errorf("success = false, want true")
	}
}

func TestItemHandler_Create_ValidationError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		ucErr  error
		wantSt int
	}{
		{"名前が空", domain.ErrInvalidName, http.StatusBadRequest},
		{"価格マイナス", domain.ErrInvalidPrice, http.StatusBadRequest},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			uc := &mockUseCase{err: tt.ucErr}
			h := handler.NewItemHandler(uc)

			rr := doRequest(t, h, http.MethodPost, "/api/items", map[string]interface{}{
				"name": "", "price": -1,
			})
			if rr.Code != tt.wantSt {
				t.Errorf("status = %d, want %d", rr.Code, tt.wantSt)
			}
		})
	}
}

func TestItemHandler_Create_InvalidBody(t *testing.T) {
	t.Parallel()
	h := handler.NewItemHandler(&mockUseCase{})

	req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

// ---- PUT /api/items/{id} ----

func TestItemHandler_Update_OK(t *testing.T) {
	t.Parallel()
	uc := &mockUseCase{updateOut: &usecase.Output{ID: 1, Name: "更新後", Price: 999}}
	h := handler.NewItemHandler(uc)

	rr := doRequest(t, h, http.MethodPut, "/api/items/1", map[string]interface{}{
		"name": "更新後", "price": 999,
	})
	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestItemHandler_Update_NotFound(t *testing.T) {
	t.Parallel()
	uc := &mockUseCase{err: domain.ErrNotFound}
	h := handler.NewItemHandler(uc)

	rr := doRequest(t, h, http.MethodPut, "/api/items/9999", map[string]interface{}{
		"name": "x",
	})
	if rr.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}

// ---- DELETE /api/items/{id} ----

func TestItemHandler_Delete_OK(t *testing.T) {
	t.Parallel()
	h := handler.NewItemHandler(&mockUseCase{})

	rr := doRequest(t, h, http.MethodDelete, "/api/items/1", nil)
	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestItemHandler_Delete_NotFound(t *testing.T) {
	t.Parallel()
	uc := &mockUseCase{err: domain.ErrNotFound}
	h := handler.NewItemHandler(uc)

	rr := doRequest(t, h, http.MethodDelete, "/api/items/9999", nil)
	if rr.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}

// ---- Method Not Allowed ----

func TestItemHandler_MethodNotAllowed(t *testing.T) {
	t.Parallel()
	h := handler.NewItemHandler(&mockUseCase{})

	rr := doRequest(t, h, http.MethodPatch, "/api/items", nil)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusMethodNotAllowed)
	}
}

// ---- CORS preflight ----

func TestItemHandler_CORS_Preflight(t *testing.T) {
	t.Parallel()
	h := handler.NewItemHandler(&mockUseCase{})

	req := httptest.NewRequest(http.MethodOptions, "/api/items", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusNoContent)
	}
	if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("missing CORS header")
	}
}
