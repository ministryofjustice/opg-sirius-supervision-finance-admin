package api

import (
	"net/http"
)

func healthCheck() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}
