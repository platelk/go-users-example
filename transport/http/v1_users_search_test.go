package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"go-users-example/domain/users"
	"go-users-example/infra/logger"
)

func TestBuilder_WithV1SearchUser(t *testing.T) {
	router := NewBuilder(logger.Logger{}, Config{}).WithV1SearchUser(func(ctx context.Context, req *users.SearchReq) (*users.SearchResp, error) {
		require.NotEmpty(t, req.IDs)
		require.Contains(t, req.IDs, "test1")
		require.Contains(t, req.IDs, "test2")
		return &users.SearchResp{}, nil
	}).router

	req := httptest.NewRequest("GET", "http://localhost/v1/users?id=test1&id=test2", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
