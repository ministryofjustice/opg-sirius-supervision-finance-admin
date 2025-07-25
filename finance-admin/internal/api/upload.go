package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/shared"
	"net/http"
)

func (c *Client) Upload(ctx context.Context, data shared.Upload) error {
	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(data)
	if err != nil {
		return err
	}

	req, err := c.newHubRequest(ctx, http.MethodPost, "/uploads", &body)

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}

	defer unchecked(resp.Body.Close)

	switch resp.StatusCode {
	case http.StatusOK:
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
