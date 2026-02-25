package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

type mockStorage struct {
	urls map[string]string
}

func newMockStorage() *mockStorage {
	return &mockStorage{urls: make(map[string]string)}
}

func (m *mockStorage) SaveUrl(_ context.Context, url string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("url cannot be empty")
	}
	short := "qwe123"
	m.urls[short] = url
	return short, nil
}

func (m *mockStorage) GetUrl(_ context.Context, short string) (string, error) {
	url, ok := m.urls[short]
	if !ok {
		return "", fmt.Errorf("not found")
	}
	return url, nil
}

func TestUrlHandler_PostUrl(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "valid url",
			body:       `{"url":"https://example.com"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "invalid json",
			body:       `{invalid}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid url",
			body:       `{"url":"not-a-url"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty url",
			body:       `{"url":""}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &UrlHandler{Storage: newMockStorage()}
			req := httptest.NewRequest(http.MethodPost, "/url", bytes.NewBufferString(tt.body))
			rec := httptest.NewRecorder()
			handler.PostUrl(rec, req)
			if rec.Code != tt.wantStatus {
				t.Errorf("PostUrl() status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestUrlHandler_PostUrl_ResponseBody(t *testing.T) {
	handler := &UrlHandler{Storage: newMockStorage()}
	body := `{"url":"https://example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/url", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	handler.PostUrl(rec, req)
	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["shortUrl"] == "" {
		t.Error("response missing shortUrl")
	}
}

func TestUrlHandler_GetUrl(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*mockStorage)
		body       string
		wantStatus int
	}{
		{
			name: "existing url",
			setup: func(m *mockStorage) {
				m.urls["qwe123"] = "https://example.com"
			},
			body:       `{"shortUrl":"qwe123"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "not found",
			setup:      func(m *mockStorage) {},
			body:       `{"shortUrl":"missing"}`,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "invalid json",
			setup:      func(m *mockStorage) {},
			body:       `{bad}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := newMockStorage()
			tt.setup(ms)
			handler := &UrlHandler{Storage: ms}
			req := httptest.NewRequest(http.MethodGet, "/url", bytes.NewBufferString(tt.body))
			rec := httptest.NewRecorder()

			handler.GetUrl(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("GetUrl() status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestUrlHandler_GetUrl_ResponseBody(t *testing.T) {
	ms := newMockStorage()
	ms.urls["qwe123"] = "https://example.com"
	handler := &UrlHandler{Storage: ms}

	req := httptest.NewRequest(http.MethodGet, "/url", bytes.NewBufferString(`{"shortUrl":"qwe123"}`))
	rec := httptest.NewRecorder()

	handler.GetUrl(rec, req)

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["url"] != "https://example.com" {
		t.Errorf("response url = %q, want %q", resp["url"], "https://example.com")
	}
}

func TestRedirectHandler_Redirect(t *testing.T) {
	tests := []struct {
		name       string
		shortUrl   string
		setup      func(*mockStorage)
		wantStatus int
	}{
		{
			name:     "existing short url",
			shortUrl: "qwe123",
			setup: func(m *mockStorage) {
				m.urls["qwe123"] = "https://example.com"
			},
			wantStatus: http.StatusFound,
		},
		{
			name:       "nonexistent short url",
			shortUrl:   "missing",
			setup:      func(m *mockStorage) {},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "empty short url",
			shortUrl:   "",
			setup:      func(m *mockStorage) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := newMockStorage()
			tt.setup(ms)
			handler := &RedirectHandler{Storage: ms}

			req := httptest.NewRequest(http.MethodGet, "/"+tt.shortUrl, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("shortUrl", tt.shortUrl)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			rec := httptest.NewRecorder()

			handler.Redirect(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("Redirect() status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestRedirectHandler_Redirect_Location(t *testing.T) {
	ms := newMockStorage()
	ms.urls["qwe123"] = "https://example.com"
	handler := &RedirectHandler{Storage: ms}
	req := httptest.NewRequest(http.MethodGet, "/qwe123", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("shortUrl", "qwe123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()
	handler.Redirect(rec, req)
	location := rec.Header().Get("Location")
	if location != "https://example.com" {
		t.Errorf("Redirect Location = %q, want %q", location, "https://example.com")
	}
}
