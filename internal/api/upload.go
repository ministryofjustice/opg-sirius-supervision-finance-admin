package api

import (
	"bytes"
	"encoding/json"
	"github.com/opg-sirius-finance-admin/internal/model"
	"net/http"
)

func (c *Client) Upload(ctx Context, data model.Upload) error {
	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(data)
	if err != nil {
		return err
	}

	req, err := c.newBackendRequest(ctx, http.MethodPost, "/uploads", &body)

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

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
			innerMap := make(map[string]string)
			innerMap[reason] = reason
			validationErrors[reason] = innerMap
		}

		return model.ValidationError{Errors: validationErrors}
	}

	return newStatusError(resp)
}
