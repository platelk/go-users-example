package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"go-users-example/domain/users"
)

// WithV1DeleteUser will add http endpoint to delete new user
func (b *Builder) WithV1DeleteUser(deleteUser users.Delete) *Builder {
	b.router.Delete("/v1/user", func(writer http.ResponseWriter, request *http.Request) {
		req, status, err := parseDeleteRequest(request)
		if err != nil {
			writer.WriteHeader(status)
			return
		}
		res, err := deleteUser(request.Context(), req)
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

func parseDeleteRequest(request *http.Request) (*users.DeleteReq, int, error) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("can't read body: %w", err)
	}
	var req users.DeleteReq
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("can't parse json body: %w", err)
	}
	return &req, 0, nil
}
