package users

import (
	"context"
	"fmt"
	"time"

	"go-users-example/infra/logger"
)

// DeleteReq contains the required parameters to delete a new user
type DeleteReq struct {
	ID string `json:"id"`
}

// DeleteResp contains the field which will be returned on successful user creation
type DeleteResp struct {
	User *User `json:"user"`
}

// Adder will save a new user in the system and generate an id for the User
type Deleter interface {
	Delete(ctx context.Context, user *User) (*User, error)
}

// Delete define the function which will delete a user in the system
type Delete func(ctx context.Context, req *DeleteReq) (*DeleteResp, error)

// SetupDelete will return a configured Delete function which can be used later
func SetupDelete(log logger.Logger, notifier ChangeNotifier, repo Deleter) Delete {
	log = log.With().Str("usecase", "user_delete").Logger()
	return notifyDelete(log, notifier, deleteUser(repo))
}

func deleteUser(repo Deleter) Delete {
	return func(ctx context.Context, req *DeleteReq) (*DeleteResp, error) {
		newUser, err := repo.Delete(ctx, &User{ID: req.ID})
		if err != nil {
			return nil, fmt.Errorf("can't delete user: %w", err)
		}
		return &DeleteResp{User: newUser}, nil
	}
}

func notifyDelete(log logger.Logger, notifier ChangeNotifier, deleteFunc Delete) Delete {
	log = log.With().Str("us_middleware", "notifier").Logger()
	return func(ctx context.Context, req *DeleteReq) (*DeleteResp, error) {
		res, err := deleteFunc(ctx, req)
		if err != nil {
			return res, err
		}
		go func(u User) {
			evt := &ChangeEvent{
				Time:   time.Now(),
				Op:     DeleteOp,
				Before: nil,
				After:  &u,
			}
			log.Debug().Interface("user", u).Msg("notify user creation")
			if err := notifier.Notify(evt); err != nil {
				log.Error().Interface("user", u).Err(err).Msg("can't send user creation event")
			}
		}(*res.User)
		return res, nil
	}
}
