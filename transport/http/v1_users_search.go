package http

import (
	"encoding/json"
	"net/http"

	"go-users-example/domain/users"
)

// WithV1SearchUser will add http endpoint to search users
// Note: here pagination is not implemented, so too many users can break the response
func (b *Builder) WithV1SearchUser(searchUser users.Search) *Builder {
	b.router.Get("/v1/users", func(writer http.ResponseWriter, request *http.Request) {
		req, status, err := parseSearchRequest(request)
		if err != nil {
			writer.WriteHeader(status)
			return
		}
		res, err := searchUser(request.Context(), req)
		switch {
		case err != nil:
			b.log.Error().Err(err).Send()
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
		default:
			data, _ := json.Marshal(res)
			_, _ = writer.Write(data)
		}
	})
	return b
}

func parseSearchRequest(request *http.Request) (*users.SearchReq, int, error) {
	return &users.SearchReq{
		IDs:       request.URL.Query()["id"],
		Emails:    request.URL.Query()["email"],
		FirstName: request.URL.Query()["first_name"],
		LastName:  request.URL.Query()["last_name"],
		NickName:  request.URL.Query()["nick_name"],
		Country:   request.URL.Query()["country"],
	}, 0, nil
}
