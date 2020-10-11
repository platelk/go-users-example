package userstore

import (
	"context"
	"fmt"

	"github.com/satori/go.uuid"

	"go-users-example/domain/users"
)

// InMemory is a user repo implementation which will store inmemory the users.
type InMemory struct {
	dataByID    map[string]*users.User
	dataEmailID map[string]string
}

// NewInMemory will initialise the store
func NewInMemory() *InMemory {
	return &InMemory{dataByID: make(map[string]*users.User), dataEmailID: make(map[string]string)}
}

// Add implements users.Adder
func (i *InMemory) Add(ctx context.Context, user *users.User) (*users.User, error) {
	if _, ok := i.dataEmailID[user.Email]; ok {
		return nil, fmt.Errorf("email %s already created: %w", user.Email, ErrAlreadyExist)
	}

	user.ID = uuid.NewV4().String()
	i.dataByID[user.ID] = &(*user) // nolint
	i.dataEmailID[user.Email] = user.ID

	return user, nil
}

// Delete will remove the user from the system
func (i *InMemory) Delete(ctx context.Context, user *users.User) (*users.User, error) {
	usr, ok := i.dataByID[user.ID]
	if !ok {
		return nil, ErrNotFound
	}

	delete(i.dataEmailID, usr.Email)
	delete(i.dataByID, usr.ID)

	return usr, nil
}

// Update will update user with same ID to the new value
func (i *InMemory) Update(ctx context.Context, user *users.User) (*users.User, error) {
	storedUser, ok := i.dataByID[user.ID]
	if !ok {
		return nil, ErrNotFound
	}

	if user.Email != "" {
		delete(i.dataEmailID, storedUser.Email)
		storedUser.Email = user.Email
		i.dataEmailID[storedUser.Email] = storedUser.ID
	}
	if user.FirstName != "" {
		storedUser.FirstName = user.FirstName
	}
	if user.LastName != "" {
		storedUser.LastName = user.LastName
	}
	if user.NickName != "" {
		storedUser.NickName = user.NickName
	}
	if user.Password != "" {
		storedUser.Password = user.Password
	}
	if user.Country != "" {
		storedUser.Country = user.Country
	}

	return storedUser, nil
}

// Query will create a query to search users. implements users.Queryier
func (i *InMemory) Query() users.Queryer {
	return &query{}
}

// Search will execute the search on the user base
func (i *InMemory) Search(ctx context.Context, q users.Queryer) ([]*users.User, error) {
	sQuery, ok := q.(*query)
	if !ok {
		return nil, ErrQueryNotCompatible
	}
	var res []*users.User
	for _, usr := range i.dataByID {
		if sQuery.match(usr) {
			res = append(res, usr)
		}
	}
	return res, nil
}

// -- internal implementation --

type query struct {
	ids       []string
	email     []string
	firstName []string
	lastName []string
	nickName []string
	country []string
}

func (q *query) ByID(id string) users.Queryer {
	q.ids = append(q.ids, id)
	return q
}

func (q *query) ByEmail(email string) users.Queryer {
	q.email = append(q.email, email)
	return q
}

func (q *query) ByFirstName(firstName string) users.Queryer {
	q.firstName = append(q.firstName, firstName)
	return q
}

func (q *query) ByLastName(lastName string) users.Queryer {
	q.lastName = append(q.lastName, lastName)
	return q
}

func (q *query) ByNickName(nickName string) users.Queryer {
	q.nickName = append(q.nickName, nickName)
	return q
}

func (q *query) ByCountry(country string) users.Queryer {
	q.country = append(q.country, country)
	return q
}

func (q *query) match(u *users.User) bool {
	for _, id := range q.ids {
		if id == u.ID {
			return true
		}
	}
	for _, email := range q.email {
		if email == u.Email {
			return true
		}
	}
	for _, firstName := range q.firstName {
		if firstName == u.FirstName {
			return true
		}
	}
	for _, lastName := range q.lastName {
		if lastName == u.LastName {
			return true
		}
	}
	for _, nickName := range q.nickName {
		if nickName == u.NickName {
			return true
		}
	}
	for _, country := range q.country {
		if country == u.Country {
			return true
		}
	}
	return false
}
