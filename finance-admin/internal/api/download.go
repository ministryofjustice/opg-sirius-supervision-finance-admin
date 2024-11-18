package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/opg-sirius-finance-admin/finance-admin/internal/model"
	"net/http"
)

func (c *Client) Download(ctx Context, uid string) (*http.Response, error) {
	req, err := c.newBackendRequest(ctx, http.MethodGet, fmt.Sprintf("/download?uid=%s", uid), nil)

	if err != nil {
		return nil, err
	}

	req, err := c.newBackendRequest(ctx, http.MethodPost, "/downloads", &body)

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println(resp.Body)

	switch resp.StatusCode {
	case http.StatusCreated:
		return nil

	case http.StatusUnauthorized:
		return ErrUnauthorized

	case http.StatusUnprocessableEntity:
		var v model.ValidationError
		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil && len(v.Errors) > 0 {
			return model.ValidationError{Errors: v.Errors}
		}

	case http.StatusBadRequest:
		var badRequests model.BadRequests
		if err := json.NewDecoder(resp.Body).Decode(&badRequests); err != nil {
			return err
		}

		validationErrors := model.ValidationErrors{}
		for _, reason := range badRequests.Reasons {
			validationErrors[reason] = map[string]string{reason: reason}
		}

		return model.ValidationError{Errors: validationErrors}
	}

	return newStatusError(resp)
}
