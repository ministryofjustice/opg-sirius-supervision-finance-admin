package api

import (
	"bytes"
	"encoding/json"
	"github.com/opg-sirius-finance-admin/internal/model"
	"net/http"
	"os"
)

func (c *Client) Upload(ctx Context, reportUploadType string, uploadDate string, email string, file *os.File) error {
	var body bytes.Buffer
	var uploadDateTransformed *model.Date
	var req *http.Request

	if uploadDate != "" {
		uploadDateFormatted := model.NewDate(uploadDate)
		uploadDateTransformed = &uploadDateFormatted
	}

	err := json.NewEncoder(&body).Encode(model.Upload{
		ReportUploadType: reportUploadType,
		UploadDate:       uploadDateTransformed,
		Email:            email,
		File:             file,
	})
	if err != nil {
		return err
	}

	switch reportUploadType {
	case "DebtChase":
		req, err = c.newSiriusRequest(ctx, http.MethodPost, "/finance/reports/upload-fee-chase", &body)
	case "DeputySchedule":
		req, err = c.newSiriusRequest(ctx, http.MethodPost, "/finance/reports/upload-deputy-billing-schedule", &body)
	default:
		req, err = c.newBackendRequest(ctx, http.MethodPost, "/uploads", &body)
	}

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
