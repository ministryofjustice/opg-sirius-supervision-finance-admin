package api

import (
	"bytes"
	"encoding/json"
	"github.com/opg-sirius-finance-admin/internal/model"
	"net/http"
	"os"
)

func (c *ApiClient) Upload(ctx Context, reportUploadType string, uploadDate string, email string, file *os.File) error {
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

	if reportUploadType == "DebtChase" {
		req, err = c.newSiriusRequest(ctx, http.MethodPost, "/finance/reports/upload-fee-chase", &body)
	} else if reportUploadType == "DeputySchedule" {
		req, err = c.newSiriusRequest(ctx, http.MethodPost, "/finance/reports/upload-deputy-billing-schedule", &body)
	} else {
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

	if resp.StatusCode == http.StatusCreated {
		return nil
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode == http.StatusUnprocessableEntity {
		var v model.ValidationError
		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil && len(v.Errors) > 0 {
			return model.ValidationError{Errors: v.Errors}
		}
	}

	if resp.StatusCode == http.StatusBadRequest {
		var badRequests model.BadRequests
		if err := json.NewDecoder(resp.Body).Decode(&badRequests); err != nil {
			return err
		}

		validationErrors := make(model.ValidationErrors)
		for _, reason := range badRequests.Reasons {
			innerMap := make(map[string]string)
			innerMap[reason] = reason
			validationErrors[reason] = innerMap
		}

		return model.ValidationError{Errors: validationErrors}
	}

	return newStatusError(resp)
}
