package users

import (
	"context"
	"fmt"
	"time"

	"go-users-example/infra/logger"
)

// UpdateReq contains the required parameters to Update a new user
type UpdateReq struct {
	ID          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	NickName    string `json:"nick_name"`
	Email       string `json:"email"`
	RawPassword string `json:"password"`
	Country     string `json:"country"`
}

// UpdateResp contains the field which will be returned on successful user update
type UpdateResp struct {
	User *User `json:"user"`
}

// Updater will update a user in the system
type Updater interface {
	Update(ctx context.Context, user *User) (*User, error)
}

// Update define the function which will Update a user in the system
type Update func(ctx context.Context, req *UpdateReq) (*UpdateResp, error)

// SetupUpdate will return a configured Update function which can be used later
func SetupUpdate(log logger.Logger, notifier ChangeNotifier, repo Updater) Update {
	log = log.With().Str("usecase", "user_update").Logger()
	return validateUpdate(log, notifyUpdate(log, notifier, updateUser(repo)))
}

func updateUser(repo Updater) Update {
	return func(ctx context.Context, req *UpdateReq) (*UpdateResp, error) {
		newUser, err := repo.Update(ctx, &User{
			ID:        req.ID,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			NickName:  req.NickName,
			Password:  req.RawPassword,
			Email:     req.Email,
			Country:   req.Country,
		})
		if err != nil {
			return nil, fmt.Errorf("can't save new user: %w", err)
		}
		return &UpdateResp{User: newUser}, nil
	}
}

func notifyUpdate(log logger.Logger, notifier ChangeNotifier, UpdateFunc Update) Update {
	log = log.With().Str("us_middleware", "notifier").Logger()
	return func(ctx context.Context, req *UpdateReq) (*UpdateResp, error) {
		res, err := UpdateFunc(ctx, req)
		if err != nil {
			return res, err
		}
		go func(u User) {
			evt := &ChangeEvent{
				Time:   time.Now(),
				Op:     UpdateOp,
				Before: nil,
				After:  &u,
			}
			log.Debug().Interface("user", u).Msg("notify user update")
			if err := notifier.Notify(evt); err != nil {
				log.Error().Interface("user", u).Err(err).Msg("can't send user update event")
			}
		}(*res.User)
		return res, nil
	}
}

func validateUpdate(log logger.Logger, updateFunc Update) Update {
	return func(ctx context.Context, req *UpdateReq) (*UpdateResp, error) {
		log.Debug().Interface("req", req).Msg("receive update")
		if err := validateFirstName(req.FirstName); req.FirstName != "" && err != nil {
			return nil, fmt.Errorf("can't validate user: %s: %w", err.Error(), ErrInvalidUser)
		}
		if err := validateLastName(req.LastName); req.LastName != "" && err != nil {
			return nil, fmt.Errorf("can't validate user: %s: %w", err.Error(), ErrInvalidUser)
		}
		if err := validateNickName(req.NickName); req.NickName != "" && err != nil {
			return nil, fmt.Errorf("can't validate user: %s: %w", err.Error(), ErrInvalidUser)
		}
		if err := validateEmail(req.Email); req.Email != "" && err != nil {
			return nil, fmt.Errorf("can't validate user: %s: %w", err.Error(), ErrInvalidUser)
		}
		return updateFunc(ctx, req)
	}
}
