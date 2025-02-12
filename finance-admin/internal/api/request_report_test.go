package api

import (
	"bytes"
	"encoding/json"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/shared"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestReport(t *testing.T) {
	mockClient := &MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000", "", "")
	dateOfTransaction := shared.NewDate("2024-05-11")
	dateTo := shared.NewDate("2025-06-15")
	dateFrom := shared.NewDate("2022-07-21")

	data := shared.ReportRequest{
		ReportType:             shared.ReportsTypeAccountsReceivable,
		JournalType:            nil,
		ScheduleType:           shared.ParseScheduleType("MOTOCardPayments"),
		AccountsReceivableType: nil,
		DebtType:               nil,
		TransactionDate:        &dateOfTransaction,
		ToDate:                 &dateTo,
		FromDate:               &dateFrom,
		Email:                  "Something@example.com",
	}

	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}, nil
	}

	err := client.RequestReport(getContext(nil), data)
	assert.NoError(t, err)
}

func TestRequestReportUnauthorised(t *testing.T) {
	mockClient := &MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000", "", "")

	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}, nil
	}

	err := client.RequestReport(getContext(nil), shared.ReportRequest{})

	assert.Equal(t, ErrUnauthorized.Error(), err.Error())
}

func TestRequestReportReturnsBadRequestError(t *testing.T) {
	mockClient := &MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000", "", "")

	json := `{"reasons":["StartDate","EndDate"]}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       r,
		}, nil
	}

	err := client.RequestReport(getContext(nil), shared.ReportRequest{})

	expectedError := model.ValidationError{Message: "", Errors: model.ValidationErrors{"EndDate": map[string]string{"EndDate": "EndDate"}, "StartDate": map[string]string{"StartDate": "StartDate"}}}
	assert.Equal(t, expectedError, err)
}

func TestRequestReportReturnsValidationError(t *testing.T) {
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

	client, _ := NewClient(http.DefaultClient, svr.URL, svr.URL, svr.URL)

	err := client.RequestReport(getContext(nil), shared.ReportRequest{})
	expectedError := model.ValidationError{Message: "", Errors: model.ValidationErrors{"ReportType": map[string]string{"required": "Please select a report type"}}}
	assert.Equal(t, expectedError, err.(model.ValidationError))
}
