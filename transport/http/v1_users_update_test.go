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

func TestBuilder_WithV1UpdateUser(t *testing.T) {
	router := NewBuilder(logger.Logger{}, Config{}).WithV1UpdateUser(func(ctx context.Context, req *users.UpdateReq) (*users.UpdateResp, error) {
		require.NotEmpty(t, req.ID)
		return &users.UpdateResp{User: &users.User{
			ID:        req.ID,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			NickName:  req.NickName,
			Password:  req.RawPassword,
			Email:     req.Email,
		}}, nil
	}).router

	req := httptest.NewRequest("PUT", "http://localhost/v1/user", strings.NewReader(`
	{"id": "testid", "first_name": "test", "last_name": "test", "email": "test@test.com", "nick_name": "test", "password": "test"}
	`))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
