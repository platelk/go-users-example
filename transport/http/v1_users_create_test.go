package http

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"go-users-example/domain/users"
	"go-users-example/infra/logger"
)

func TestBuilder_WithV1CreateUser(t *testing.T) {
	router := NewBuilder(logger.Logger{}, Config{}).WithV1CreateUser(func(ctx context.Context, req *users.CreateReq) (*users.CreateResp, error) {
		require.NotEmpty(t, req.Email)
		require.NotEmpty(t, req.NickName)
		require.NotEmpty(t, req.FirstName)
		require.NotEmpty(t, req.LastName)
		return &users.CreateResp{User: &users.User{
			ID:        "newid",
			FirstName: req.FirstName,
			LastName:  req.LastName,
			NickName:  req.NickName,
			Password:  req.RawPassword,
			Email:     req.Email,
		}}, nil
	}).router

	req := httptest.NewRequest("POST", "http://localhost/v1/user", strings.NewReader(`
	{"first_name": "test", "last_name": "test", "email": "test@test.com", "nick_name": "test", "password": "test"}
	`))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Contains(t, string(body), "newid")
}
