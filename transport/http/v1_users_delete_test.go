package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"go-users-example/domain/users"
	"go-users-example/infra/logger"
)

func TestBuilder_WithV1DeleteUser(t *testing.T) {
	router := NewBuilder(logger.Logger{}, Config{}).WithV1DeleteUser(func(ctx context.Context, req *users.DeleteReq) (*users.DeleteResp, error) {
		require.NotEmpty(t, req.ID)
		return &users.DeleteResp{User: &users.User{
			ID: req.ID,
		}}, nil
	}).router

	req := httptest.NewRequest("DELETE", "http://localhost/v1/user", strings.NewReader(`
	{"id": "testid"}
	`))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
