package http

import (
	"net/http"
)

// WithHealthCheck will add default endpoints to check its status
func (b *Builder) WithHealthCheck() *Builder {
	b.router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})
	b.router.Get("/ping", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("pong"))
	})
	return b
}
