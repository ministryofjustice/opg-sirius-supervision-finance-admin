package api

import (
	"bytes"
	"encoding/json"
	"github.com/opg-sirius-finance-admin/finance-admin/internal/model"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSubmitDownload(t *testing.T) {
	mockClient := &MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000", "")
	dateOfTransaction := model.NewDate("2024-05-11")
	dateTo := model.NewDate("2025-06-15")
	dateFrom := model.NewDate("2022-07-21")

	data := model.Download{
		ReportType:         "reportType",
		ReportJournalType:  "reportJournalType",
		ReportScheduleType: "reportScheduleType",
		ReportAccountType:  "reportAccountType",
		ReportDebtType:     "reportDebtType",
		DateOfTransaction:  &dateOfTransaction,
		ToDateField:        &dateTo,
		FromDateField:      &dateFrom,
		Email:              "Something@example.com",
	}

	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}, nil
	}

	err := client.Download(getContext(nil), data)
	assert.NoError(t, err)
}

func TestSubmitDownloadUnauthorised(t *testing.T) {
	mockClient := &MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000", "")

	data := model.Download{
		ReportType:         "reportType",
		ReportJournalType:  "reportJournalType",
		ReportScheduleType: "reportScheduleType",
		ReportAccountType:  "reportAccountType",
		ReportDebtType:     "reportDebtType",
		DateOfTransaction:  nil,
		ToDateField:        nil,
		FromDateField:      nil,
		Email:              "Something@example.com",
	}

	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}, nil
	}

	err := client.Download(getContext(nil), data)

	assert.Equal(t, ErrUnauthorized.Error(), err.Error())
}

func TestSubmitDownloadReturnsBadRequestError(t *testing.T) {
	mockClient := &MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000", "")

	data := model.Download{
		ReportType:         "reportType",
		ReportJournalType:  "reportJournalType",
		ReportScheduleType: "reportScheduleType",
		ReportAccountType:  "reportAccountType",
		ReportDebtType:     "reportDebtType",
		DateOfTransaction:  nil,
		ToDateField:        nil,
		FromDateField:      nil,
		Email:              "Something@example.com",
	}

	json := `{"reasons":["StartDate","EndDate"]}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       r,
		}, nil
	}

	err := client.Download(getContext(nil), data)

	expectedError := model.ValidationError{Message: "", Errors: model.ValidationErrors{"EndDate": map[string]string{"EndDate": "EndDate"}, "StartDate": map[string]string{"StartDate": "StartDate"}}}
	assert.Equal(t, expectedError, err)
}

func TestSubmitDownloadReturnsValidationError(t *testing.T) {
	data := model.Download{
		ReportType:         "",
		ReportJournalType:  "reportJournalType",
		ReportScheduleType: "reportScheduleType",
		ReportAccountType:  "reportAccountType",
		ReportDebtType:     "reportDebtType",
		DateOfTransaction:  nil,
		ToDateField:        nil,
		FromDateField:      nil,
		Email:              "Something@example.com",
	}

	validationErrors := model.ValidationError{
		Message: "Validation failed",
		Errors: map[string]map[string]string{
			"ReportType": {
				"required": "Please select a report type",
			},
		},
	}
	responseBody, _ := json.Marshal(validationErrors)
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write(responseBody)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL, svr.URL)

	err := client.Download(getContext(nil), data)
	expectedError := model.ValidationError{Message: "", Errors: model.ValidationErrors{"ReportType": map[string]string{"required": "Please select a report type"}}}
	assert.Equal(t, expectedError, err.(model.ValidationError))
}
