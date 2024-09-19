package api

import (
	"fmt"
	"net/http"
)

func (s *Server) upload(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("Hello :)")
	return nil
}
