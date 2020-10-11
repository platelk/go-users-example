package users

import (
	"context"
	"fmt"

	"go-users-example/infra/logger"
)

// SearchReq contains the required parameters to search users
type SearchReq struct {
	IDs       []string
	Emails    []string
	FirstName []string
	LastName []string
	NickName []string
	Country []string
}

// SearchResp contains the field which will be returned on successful user search
type SearchResp struct {
	Users []*User `json:"users"`
}

// Queryer will construct the query to search users based on multiple criteria
type Queryer interface {
	ByID(id string) Queryer
	ByEmail(email string) Queryer
	ByFirstName(firstName string) Queryer
	ByLastName(lastName string) Queryer
	ByNickName(nickName string) Queryer
	ByCountry(country string) Queryer
}

// Searcher will allow searching users based on different criteria
type Searcher interface {
	Query() Queryer
	Search(ctx context.Context, query Queryer) ([]*User, error)
}

// Search define the function which will search for users in the system
type Search func(ctx context.Context, req *SearchReq) (*SearchResp, error)

// SetupSearch will return a configured Create function which can be used later
func SetupSearch(log logger.Logger, repo Searcher) Search {
	log = log.With().Str("usecase", "user_search").Logger()
	return searchUser(repo)
}

func searchUser(repo Searcher) Search {
	return func(ctx context.Context, req *SearchReq) (*SearchResp, error) {
		qBuilder := repo.Query()

		for _, id := range req.IDs {
			qBuilder = qBuilder.ByID(id)
		}
		for _, email := range req.Emails {
			qBuilder = qBuilder.ByEmail(email)
		}
		for _, firstName := range req.FirstName {
			qBuilder = qBuilder.ByFirstName(firstName)
		}
		for _, lastName := range req.LastName {
			qBuilder = qBuilder.ByLastName(lastName)
		}
		for _, nickName := range req.NickName {
			qBuilder = qBuilder.ByNickName(nickName)
		}
		for _, country := range req.Country {
			qBuilder = qBuilder.ByCountry(country)
		}

		users, err := repo.Search(ctx, qBuilder)
		if err != nil {
			return nil, fmt.Errorf("can't perform search: %w", err)
		}

		return &SearchResp{Users: users}, nil
	}
}
