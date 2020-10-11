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

func TestSetupUpdate_OK(t *testing.T) {
	userStore := userstore.NewInMemory()
	usr, _ := userStore.Add(context.Background(), &users.User{
		Email: "test-update-1",
	})
	update := users.SetupUpdate(logger.Logger{}, usernotifier.NewInMemory(), userStore)
	res, err := update(context.Background(), &users.UpdateReq{
		ID:    usr.ID,
		Email: "test-update-1-updated@test.com",
	})
	require.NoError(t, err)
	require.Equal(t, res.User.Email, "test-update-1-updated@test.com")
	require.NotEmpty(t, res.User.ID)
}
