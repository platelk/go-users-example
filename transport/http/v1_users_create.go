package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"go-users-example/domain/users"
)

// WithV1CreateUser will add http endpoint to create new user
func (b *Builder) WithV1CreateUser(createUser users.Create) *Builder {
	b.router.Post("/v1/user", func(writer http.ResponseWriter, request *http.Request) {
		req, status, err := parseRequest(request)
		if err != nil {
			writer.WriteHeader(status)
			return
		}
		res, err := createUser(request.Context(), req)
		switch {
		case errors.Is(err, users.ErrInvalidUser):
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte(err.Error()))
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

func parseRequest(request *http.Request) (*users.CreateReq, int, error) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("can't read body: %w", err)
	}
	var req users.CreateReq
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("can't parse json body: %w", err)
	}

	return &req, 0, nil
}
