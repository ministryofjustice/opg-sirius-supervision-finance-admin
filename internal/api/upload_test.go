package api

import (
	"bytes"
	"encoding/json"
	"github.com/opg-sirius-finance-admin/internal/model"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUploadUrlSwitching(t *testing.T) {
	mockClient := &MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000", "")
	uploadDate := model.Date{Time: time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)}
	content := []byte("file content")

	data := model.Upload{
		ReportUploadType: "reportUploadType",
		UploadDate:       &uploadDate,
		Email:            "Something@example.com",
		File:             content,
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

	client, _ := NewClient(http.DefaultClient, svr.URL, svr.URL)

	err := client.Upload(getContext(nil), model.Upload{})

	assert.Equal(t, ErrUnauthorized.Error(), err.Error())
}

func TestSubmitUploadReturns500Error(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL, svr.URL)

	err := client.Upload(getContext(nil), model.Upload{})

	assert.Equal(t, StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/uploads",
		Method: http.MethodPost,
	}, err)
}

func TestSubmitUploadReturnsBadRequestError(t *testing.T) {
	mockClient := &MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000", "")

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

	err := client.Upload(getContext(nil), model.Upload{})

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

	client, _ := NewClient(http.DefaultClient, svr.URL, svr.URL)

	err := client.Upload(getContext(nil), model.Upload{})
	expectedError := model.ValidationError{Message: "", Errors: model.ValidationErrors{"ReportUploadType": map[string]string{"required": "Please select a report type"}}}
	assert.Equal(t, expectedError, err.(model.ValidationError))
}
