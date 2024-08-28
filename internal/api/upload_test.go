package api

import (
	"bytes"
	"encoding/json"
	"github.com/opg-sirius-finance-admin/internal/model"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadUrlSwitching(t *testing.T) {
	mockClient := &MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000", "")

	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name()) // Clean up after the test
	defer tempFile.Close()

	// Write some content to the temp file
	content := []byte("fake file content")
	if _, err := tempFile.Write(content); err != nil {
		t.Fatal(err)
	}
	// Reset the file pointer to the beginning
	if _, err := tempFile.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	// Define test cases
	testCases := map[string]string{
		"DebtChase":      "http://localhost:3000/supervision-api/v1/finance/reports/upload-fee-chase",
		"DeputySchedule": "http://localhost:3000/supervision-api/v1/finance/reports/upload-deputy-billing-schedule",
		"OtherType":      "/uploads",
	}

	for reportUploadType, expectedURL := range testCases {
		t.Run(reportUploadType, func(t *testing.T) {
			// Variable to capture the request URL
			var capturedURL *url.URL

			// Mock the HTTP client's Do function to capture the request URL
			GetDoFunc = func(req *http.Request) (*http.Response, error) {
				capturedURL = req.URL
				return &http.Response{
					StatusCode: http.StatusCreated,
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				}, nil
			}

			err := client.Upload(getContext(nil), reportUploadType, "", "", tempFile)
			assert.NoError(t, err)
			assert.Equal(t, expectedURL, capturedURL.String())
		})
	}
}

func TestSubmitUploadUnauthorised(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL, svr.URL)

	err := client.Upload(getContext(nil), "", "", "", nil)

	assert.Equal(t, ErrUnauthorized.Error(), err.Error())
}

func TestSubmitUploadReturns500Error(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL, svr.URL)

	err := client.Upload(getContext(nil), "", "", "", nil)

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

	err := client.Upload(getContext(nil), "", "", "", nil)

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

	err := client.Upload(getContext(nil), "", "", "", nil)
	expectedError := model.ValidationError{Message: "", Errors: model.ValidationErrors{"ReportUploadType": map[string]string{"required": "Please select a report type"}}}
	assert.Equal(t, expectedError, err.(model.ValidationError))
}
