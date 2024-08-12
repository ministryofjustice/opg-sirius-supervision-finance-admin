package api

import (
	"bytes"
	"encoding/json"
	"github.com/opg-sirius-finance-admin/internal/model"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubmitDownload(t *testing.T) {
	mockClient := SetUpTest()
	client, _ := NewApiClient(mockClient, "http://localhost:3000", "")

	data := `{
		"reportType":         "AccountsReceivable",
		"reportJournalType":  "",
		"reportScheduleType": "",
		"reportAccountType":  "BadDebtWriteOffReport",
		"reportDebtType":     "",
		"dateField":          "11/05/2024",
		"dateFromField":      "01/04/2024",
		"dateToField":        "31/03/2025",
		"emailField":         "SomeSortOfEmail@example.com",
	}
	`

	r := io.NopCloser(bytes.NewReader([]byte(data)))

	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 201,
			Body:       r,
		}, nil
	}

	err := client.Download(getContext(nil), "AccountsReceivable", "", "", "BadDebtWriteOffReport", "", "11/05/2024", "01/04/2024", "31/03/2025", "SomeSortOfEmail@example.com")
	assert.Equal(t, nil, err)
}

func TestSubmitDownloadUnauthorised(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, svr.URL)

	err := client.Download(getContext(nil), "AccountsReceivable", "", "", "BadDebtWriteOffReport", "", "11/05/2024", "01/04/2024", "31/03/2025", "SomeSortOfEmail@example.com")

	assert.Equal(t, ErrUnauthorized.Error(), err.Error())
}

func TestSubmitDownloadReturns500Error(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, svr.URL)

	err := client.Download(getContext(nil), "AccountsReceivable", "", "", "BadDebtWriteOffReport", "", "11/05/2024", "01/04/2024", "31/03/2025", "SomeSortOfEmail@example.com")

	assert.Equal(t, StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/downloads",
		Method: http.MethodGet,
	}, err)
}

func TestSubmitDownloadReturnsBadRequestError(t *testing.T) {
	mockClient := SetUpTest()
	client, _ := NewApiClient(mockClient, "http://localhost:3000", "")

	json := `
		{"reasons":["StartDate","EndDate"]}
	`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 400,
			Body:       r,
		}, nil
	}

	err := client.Download(getContext(nil), "AccountsReceivable", "", "", "BadDebtWriteOffReport", "", "11/05/2024", "01/04/2024", "31/03/2025", "SomeSortOfEmail@example.com")

	expectedError := model.ValidationError{Message: "", Errors: model.ValidationErrors{"EndDate": map[string]string{"EndDate": "EndDate"}, "StartDate": map[string]string{"StartDate": "StartDate"}}}
	assert.Equal(t, expectedError, err)
}

func TestSubmitDownloadReturnsValidationError(t *testing.T) {
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

	client, _ := NewApiClient(http.DefaultClient, svr.URL, svr.URL)

	err := client.Download(getContext(nil), "", "", "", "", "", "", "", "", "")
	expectedError := model.ValidationError{Message: "", Errors: model.ValidationErrors{"ReportType": map[string]string{"required": "Please select a report type"}}}
	assert.Equal(t, expectedError, err.(model.ValidationError))
}
