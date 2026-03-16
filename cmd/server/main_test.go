package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestSetupServer はサーバーの初期設定をテストします。
func TestSetupServer(t *testing.T) {
	e := setupServer()
	if e == nil {
		t.Fatal("setupServer()がnilを返しました")
	}
}

// TestHealthEndpoint はヘルスチェックエンドポイントをテストします。
func TestHealthEndpoint(t *testing.T) {
	e := setupServer()

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{
			name:       "ヘルスチェック正常",
			method:     http.MethodGet,
			path:       "/health",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("ステータスコード = %d, 期待値 %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
