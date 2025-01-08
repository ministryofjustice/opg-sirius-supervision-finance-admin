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

func TestUploadUrlSwitching(t *testing.T) {
	mockClient := &MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000", "", "")

	data := shared.Upload{
		ReportUploadType: shared.ParseReportUploadType("reportUploadType"),
		UploadDate:       shared.NewDate("2025-06-15"),
		Email:            "Something@example.com",
		File:             []byte("file content"),
	}

	GetDoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}, nil
	}

	err := client.Upload(getContext(nil), data)
	assert.NoError(t, err)
}

func TestSubmitUploadUnauthorised(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL, svr.URL, "")

	err := client.Upload(getContext(nil), shared.Upload{})

	assert.Equal(t, ErrUnauthorized.Error(), err.Error())
}

func TestSubmitUploadReturns500Error(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL, svr.URL, "")

	err := client.Upload(getContext(nil), shared.Upload{})

	assert.Equal(t, StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/uploads",
		Method: http.MethodPost,
	}, err)
}

func TestSubmitUploadReturnsBadRequestError(t *testing.T) {
	mockClient := &MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000", "", "")

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

	err := client.Upload(getContext(nil), shared.Upload{})

	expectedError := model.ValidationError{Message: "", Errors: model.ValidationErrors{"EndDate": map[string]string{"EndDate": "EndDate"}, "StartDate": map[string]string{"StartDate": "StartDate"}}}
	assert.Equal(t, expectedError, err)
}

func TestSubmitUploadReturnsValidationError(t *testing.T) {
	validationErrors := model.ValidationError{
		Message: "Validation failed",
		Errors: map[string]map[string]string{
			"ReportUploadType": {
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

	client, _ := NewClient(http.DefaultClient, svr.URL, svr.URL, "")

	err := client.Upload(getContext(nil), shared.Upload{})
	expectedError := model.ValidationError{Message: "", Errors: model.ValidationErrors{"ReportUploadType": map[string]string{"required": "Please select a report type"}}}
	assert.Equal(t, expectedError, err.(model.ValidationError))
}
