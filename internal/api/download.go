package api

import (
	"bytes"
	"encoding/json"
	"github.com/opg-sirius-finance-admin/internal/model"
	"net/http"
)

func (c *Client) Download(ctx Context, reportType string, reportJournalType string, reportScheduleType string, reportAccountType string, reportDebtType string, dateOfTransaction string, dateFrom string, dateTo string, email string) error {
	var body bytes.Buffer
	dateTransformed, toDateTransformed, fromDateTransformed := NewDownload(dateOfTransaction, dateTo, dateFrom)

	err := json.NewEncoder(&body).Encode(model.Download{
		ReportType:         reportType,
		ReportJournalType:  reportJournalType,
		ReportScheduleType: reportScheduleType,
		ReportAccountType:  reportAccountType,
		ReportDebtType:     reportDebtType,
		DateOfTransaction:  dateTransformed,
		ToDateField:        toDateTransformed,
		FromDateField:      fromDateTransformed,
		Email:              email,
	})
	if err != nil {
		return err
	}

	req, err := c.newBackendRequest(ctx, http.MethodGet, "/downloads", &body)

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

func NewDownload(dateOfTransaction string, dateTo string, dateFrom string) (*model.Date, *model.Date, *model.Date) {
	var dateTransformed *model.Date
	var toDateTransformed *model.Date
	var fromDateTransformed *model.Date

	if dateOfTransaction != "" {
		raisedDateFormatted := model.NewDate(dateOfTransaction)
		dateTransformed = &raisedDateFormatted
	}

	if dateTo != "" {
		startDateFormatted := model.NewDate(dateTo)
		toDateTransformed = &startDateFormatted
	}

	if dateFrom != "" {
		endDateFormatted := model.NewDate(dateFrom)
		fromDateTransformed = &endDateFormatted
	}
	return dateTransformed, toDateTransformed, fromDateTransformed
}
