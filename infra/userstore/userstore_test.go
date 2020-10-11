package userstore

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"go-users-example/domain/users"
)

type wrongQuery struct {}

func (w *wrongQuery) ByLastName(lastName string) users.Queryer {
	panic("implement me")
}

func (w *wrongQuery) ByNickName(nickName string) users.Queryer {
	panic("implement me")
}

func (w *wrongQuery) ByCountry(country string) users.Queryer {
	panic("implement me")
}

func (w *wrongQuery) ByID(id string) users.Queryer {
	panic("implement me")
}

func (w *wrongQuery) ByEmail(email string) users.Queryer {
	panic("implement me")
}

func (w *wrongQuery) ByFirstName(firstName string) users.Queryer {
	panic("implement me")
}

type userStore interface {
	users.Adder
	users.Updater
	users.Deleter
	users.Searcher
}

func runTestSuite(t *testing.T, store userStore) {
	runTestAdd(t, store)
	runTestDelete(t, store)
	runTestUpdate(t, store)
	runTestSearch(t, store)
}

func runTestSearch(t *testing.T, store userStore) {
	t.Run("search by ID", func(t *testing.T) {
		usr, _ := store.Add(context.Background(), &users.User{
			Email: "test-search-1",
		})
		res, _ := store.Search(context.Background(), store.Query().ByID(usr.ID))
		require.NotEmpty(t, res)
		require.Equal(t, usr.ID, res[0].ID)
	})
	t.Run("search by firstname", func(t *testing.T) {
		usr, _ := store.Add(context.Background(), &users.User{
			FirstName: "test1",
			Email:     "test-search-2",
		})
		res, _ := store.Search(context.Background(), store.Query().ByFirstName(usr.FirstName))
		require.NotEmpty(t, res)
		require.Equal(t, usr.ID, res[0].ID)
	})
	t.Run("search by email", func(t *testing.T) {
		usr, _ := store.Add(context.Background(), &users.User{
			Email: "test-search-3",
		})
		res, _ := store.Search(context.Background(), store.Query().ByEmail(usr.Email))
		require.NotEmpty(t, res)
		require.Equal(t, usr.ID, res[0].ID)
	})
	t.Run("error on no compatible query", func(t *testing.T) {
		_, err := store.Search(context.Background(), &wrongQuery{})
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrQueryNotCompatible))
	})
}

func runTestUpdate(t *testing.T, store userStore) {
	t.Run("update not found user", func(t *testing.T) {
		_, err := store.Update(context.Background(), &users.User{
			ID: "unknown",
		})
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrNotFound))
	})
	t.Run("update email", func(t *testing.T) {
		usr, err := store.Add(context.Background(), &users.User{
			FirstName: "test",
			LastName:  "test",
			NickName:  "test",
			Password:  "test",
			Email:     "test-update-1",
		})
		require.NoError(t, err)
		_, err = store.Update(context.Background(), &users.User{
			ID:        usr.ID,
			FirstName: "updated",
			LastName:  "updated",
			NickName:  "updated",
			Password:  "updated",
			Email:     "test-update-1-updated",
		})
		require.NoError(t, err)

		res, _ := store.Search(context.Background(), store.Query().ByEmail("test-update-1-updated"))
		require.NotEmpty(t, res)
		usrUpdated := res[0]

		require.Equal(t, "test-update-1-updated", usrUpdated.Email)
		require.Equal(t, "updated", usrUpdated.FirstName)
		require.Equal(t, "updated", usrUpdated.LastName)
		require.Equal(t, "updated", usrUpdated.NickName)
		require.Equal(t, "updated", usrUpdated.Password)
	})
}

func runTestDelete(t *testing.T, store userStore) {
	t.Run("delete after add", func(t *testing.T) {
		usr, err := store.Add(context.Background(), &users.User{
			FirstName: "test",
			LastName:  "test",
			NickName:  "test",
			Password:  "test",
			Email:     "test-delete-1",
		})
		require.NoError(t, err)

		res, _ := store.Search(context.Background(), store.Query().ByEmail(usr.Email))
		require.NotEmpty(t, res)

		_, err = store.Delete(context.Background(), usr)
		require.NoError(t, err)

		res, _ = store.Search(context.Background(), store.Query().ByEmail(usr.Email))
		require.Empty(t, res)
	})
	t.Run("delete unknown user", func(t *testing.T) {
		_, err := store.Delete(context.Background(), &users.User{
			FirstName: "test",
			LastName:  "test",
			NickName:  "test",
			Password:  "test",
			Email:     "test-delete-unknown",
		})
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrNotFound))
	})
}

func runTestAdd(t *testing.T, store userStore) {
	t.Run("add basic user", func(t *testing.T) {
		usr, err := store.Add(context.Background(), &users.User{
			FirstName: "test",
			LastName:  "test",
			NickName:  "test",
			Password:  "test",
			Email:     "test-add-1",
		})
		require.NoError(t, err)
		require.NotEmpty(t, usr)
	})
	t.Run("add multiple user", func(t *testing.T) {
		usr, err := store.Add(context.Background(), &users.User{
			FirstName: "test",
			LastName:  "test",
			NickName:  "test",
			Password:  "test",
			Email:     "test-add-2",
		})
		require.NoError(t, err)
		require.NotEmpty(t, usr)
		res, _ := store.Search(context.Background(), store.Query().ByEmail("test-add-1"))
		require.NotEmpty(t, res)
		require.NotEqual(t, res[0].ID, usr.ID)
	})
	t.Run("add user with same email", func(t *testing.T) {
		_, err := store.Add(context.Background(), &users.User{
			Email: "test-add-1",
		})
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrAlreadyExist))
	})
}
