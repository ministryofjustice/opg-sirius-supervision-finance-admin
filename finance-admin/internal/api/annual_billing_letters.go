package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/shared"
)

func (c *Client) AnnualBillingLetters(ctx context.Context) (shared.AnnualBillingInformation, error) {

	abi := shared.AnnualBillingInformation{}
	req, err := c.newHubRequest(ctx, http.MethodGet, "/annual-billing-letters", nil)

	if err != nil {
		return abi, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return abi, err
	}

	defer unchecked(resp.Body.Close)

	switch resp.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(resp.Body).Decode(&abi); err != nil {
			return shared.AnnualBillingInformation{}, err
		}
		return abi, nil

	case http.StatusUnauthorized:
		return abi, ErrUnauthorized

	case http.StatusUnprocessableEntity:
		var v model.ValidationError
		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil && len(v.Errors) > 0 {
			return abi, model.ValidationError{Errors: v.Errors}
		}

	case http.StatusBadRequest:
		var badRequests model.BadRequests
		if err := json.NewDecoder(resp.Body).Decode(&badRequests); err != nil {
			return abi, err
		}

		validationErrors := model.ValidationErrors{}
		for _, reason := range badRequests.Reasons {
			validationErrors[reason] = map[string]string{reason: reason}
		}

		return abi, model.ValidationError{Errors: validationErrors}
	default:
		return abi, newStatusError(resp)
	}
	return abi, newStatusError(resp)
}
