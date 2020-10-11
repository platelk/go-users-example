package users_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"go-users-example/domain/users"
	"go-users-example/infra/logger"
	"go-users-example/infra/pwdhasher"
	"go-users-example/infra/usernotifier"
	"go-users-example/infra/userstore"
)

func TestSetupCreate_OK(t *testing.T) {
	create := users.SetupCreate(logger.Logger{}, usernotifier.NewInMemory(), userstore.NewInMemory(), pwdhasher.NewBcrypt())
	res, err := create(context.Background(), &users.CreateReq{
		FirstName:   "test",
		LastName:    "test",
		NickName:    "test",
		Email:       "test-create-1@test.com",
		RawPassword: "",
	})
	require.NoError(t, err)
	require.Equal(t, res.User.Email, "test-create-1@test.com")
	require.NotEmpty(t, res.User.ID)
}
