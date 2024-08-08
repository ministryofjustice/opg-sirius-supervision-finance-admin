package api

import (
	"bytes"
	"encoding/json"
	"github.com/opg-sirius-finance-admin/internal/model"
	"net/http"
)

func (c *ApiClient) SubmitDownload(ctx Context, reportType string, reportJournalType string, reportScheduleType string, reportAccountType string, reportDebtType string, dateField string, dateFromField string, dateToField string, emailField string) error {
	var body bytes.Buffer
	var dateTransformed *model.Date
	var toDateTransformed *model.Date
	var fromDateTransformed *model.Date

	if dateField != "" {
		raisedDateFormatted := model.NewDate(dateField)
		dateTransformed = &raisedDateFormatted
	}

	if dateToField != "" {
		startDateFormatted := model.NewDate(dateToField)
		toDateTransformed = &startDateFormatted
	}

	if dateFromField != "" {
		endDateFormatted := model.NewDate(dateFromField)
		fromDateTransformed = &endDateFormatted
	}

	err := json.NewEncoder(&body).Encode(model.Download{
		ReportType:         reportType,
		ReportJournalType:  reportJournalType,
		ReportScheduleType: reportScheduleType,
		ReportAccountType:  reportAccountType,
		ReportDebtType:     reportDebtType,
		DateField:          dateTransformed,
		ToDateField:        toDateTransformed,
		FromDateField:      fromDateTransformed,
		Email:              emailField,
	})
	if err != nil {
		return err
	}

	req, err := c.newSiriusRequest(ctx, http.MethodPost, "/downloads", &body)

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
