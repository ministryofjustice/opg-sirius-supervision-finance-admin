package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/shared"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUploadSuccess(t *testing.T) {
	mockClient := &MockClient{}
	mockJwtClient := &mockJWTClient{}
	client := NewClient(mockClient, mockJwtClient, EnvVars{"http://localhost:3000", ""})

	data := shared.Upload{
		UploadType:   shared.ParseReportUploadType("reportUploadType"),
		UploadDate:   shared.NewDate("2025-06-15"),
		EmailAddress: "Something@example.com",
		Base64Data:   base64.StdEncoding.EncodeToString([]byte("col1, col2\nabc,1")),
	}

	GetDoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}, nil
	}

	err := client.Upload(testContext(), data)
	assert.NoError(t, err)
}

func TestSubmitUploadUnauthorised(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	mockJwtClient := &mockJWTClient{}
	client := NewClient(http.DefaultClient, mockJwtClient, EnvVars{svr.URL, svr.URL})

	err := client.Upload(testContext(), shared.Upload{})

	assert.Equal(t, ErrUnauthorized.Error(), err.Error())
}

func TestSubmitUploadReturns500Error(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	mockJwtClient := &mockJWTClient{}
	client := NewClient(http.DefaultClient, mockJwtClient, EnvVars{svr.URL, svr.URL})

	err := client.Upload(testContext(), shared.Upload{})

	assert.Equal(t, StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/uploads",
		Method: http.MethodPost,
	}, err)
}

func TestSubmitUploadReturnsBadRequestError(t *testing.T) {
	mockClient := &MockClient{}
	mockJwtClient := &mockJWTClient{}
	client := NewClient(mockClient, mockJwtClient, EnvVars{"http://localhost:3000", ""})

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

	err := client.Upload(testContext(), shared.Upload{})

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

	mockJwtClient := &mockJWTClient{}
	client := NewClient(http.DefaultClient, mockJwtClient, EnvVars{svr.URL, svr.URL})

	err := client.Upload(testContext(), shared.Upload{})
	expectedError := model.ValidationError{Message: "", Errors: model.ValidationErrors{"ReportUploadType": map[string]string{"required": "Please select a report type"}}}
	assert.Equal(t, expectedError, err.(model.ValidationError))
}
