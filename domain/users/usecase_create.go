package users

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-users-example/infra/logger"
)

// ErrInvalidUser is returned if one of the provided field isn't validated
var ErrInvalidUser = errors.New("provided user isn't valid")

// CreateReq contains the required parameters to create a new user
type CreateReq struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	NickName    string `json:"nick_name"`
	Email       string `json:"email"`
	Country     string `json:"country"`
	RawPassword string `json:"password"`
}

// CreateResp contains the field which will be returned on successful user creation
type CreateResp struct {
	User *User `json:"user"`
}

// Adder will save a new user in the system and generate an id for the User
type Adder interface {
	Add(ctx context.Context, user *User) (*User, error)
}

// Hasher will hash the password to securely store it
type Hasher interface {
	Hash(pwd string) (string, error)
}

// ChangeNotifier will propagate change event about the user.
// Note: Here the change is asynchronous and considered "losable".
// 		 for a more event driven architecture,
//		 the create usecase should notify only and the event should trigger the creation
type ChangeNotifier interface {
	Notify(event *ChangeEvent) error
}

// Create define the function which will create a user in the system
type Create func(ctx context.Context, req *CreateReq) (*CreateResp, error)

// SetupCreate will return a configured Create function which can be used later
func SetupCreate(log logger.Logger, notifier ChangeNotifier, repo Adder, hasher Hasher) Create {
	log = log.With().Str("usecase", "user_create").Logger()
	return validateCreate(notifyCreate(log, notifier, createUser(repo, hasher)))
}

func createUser(repo Adder, hash Hasher) Create {
	return func(ctx context.Context, req *CreateReq) (*CreateResp, error) {
		hashedPwd, err := hash.Hash(req.RawPassword)
		if err != nil {
			return nil, fmt.Errorf("can't hash the password: %w", err)
		}
		newUser, err := repo.Add(ctx, &User{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			NickName:  req.NickName,
			Password:  hashedPwd,
			Email:     req.Email,
			Country:   req.Country,
		})
		if err != nil {
			return nil, fmt.Errorf("can't save new user: %w", err)
		}
		return &CreateResp{User: newUser}, nil
	}
}

func validateCreate(createFunc Create) Create {
	return func(ctx context.Context, req *CreateReq) (*CreateResp, error) {
		err := validateUser(&User{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			NickName:  req.NickName,
			Email:     req.Email,
			Country:   req.Country,
		})
		if err != nil {
			return nil, fmt.Errorf("can't validate user: %s: %w", err.Error(), ErrInvalidUser)
		}
		return createFunc(ctx, req)
	}
}

func notifyCreate(log logger.Logger, notifier ChangeNotifier, createFunc Create) Create {
	log = log.With().Str("us_middleware", "notifier").Logger()
	return func(ctx context.Context, req *CreateReq) (*CreateResp, error) {
		res, err := createFunc(ctx, req)
		if err != nil {
			return res, err
		}
		go func(u User) {
			evt := &ChangeEvent{
				Time:   time.Now(),
				Op:     CreateOp,
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
