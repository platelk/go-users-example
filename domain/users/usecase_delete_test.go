package users_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"go-users-example/domain/users"
	"go-users-example/infra/logger"
	"go-users-example/infra/usernotifier"
	"go-users-example/infra/userstore"
)

func TestSetupDelete_OK(t *testing.T) {
	userStore := userstore.NewInMemory()
	usr, _ := userStore.Add(context.Background(), &users.User{
		Email:     "test-delete-1",
	})
	del := users.SetupDelete(logger.Logger{}, usernotifier.NewInMemory(), userStore)
	res, err := del(context.Background(), &users.DeleteReq{
		ID: usr.ID,
	})
	require.NoError(t, err)
	require.Equal(t, res.User.Email, "test-delete-1")
	require.NotEmpty(t, res.User.ID)
}
