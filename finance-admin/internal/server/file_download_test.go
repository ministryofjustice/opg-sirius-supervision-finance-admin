package server

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestDownload_success(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "download/request?uid=dGVzdC5jc3Y=", nil)
	r.Header.Add("HX-Request", "true")

	envVars := EnvironmentVars{
		SiriusURL: "https://sirius.gov.uk/v1/api",
		Prefix:    "finance-admin",
	}

	handler := requestDownload(envVars)
	handler.ServeHTTP(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "finance-admin/download/callback?uid=dGVzdC5jc3Y=", resp.Header.Get("HX-Redirect"))
}

func TestRequestDownload_notHtmx(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "download/request?uid=dGVzdC5jc3Y=", nil)

	envVars := EnvironmentVars{
		SiriusURL: "https://sirius.gov.uk/v1/api",
		Prefix:    "finance-admin",
	}

	handler := requestDownload(envVars)
	handler.ServeHTTP(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestDownloadCallback_success(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "download/callback?uid=dGVzdC5jc3Y=", nil)
	fileContent := "col1,col2,col3\n1,a,Z\n"

	download := &http.Response{
		Body: io.NopCloser(bytes.NewBufferString(fileContent)),
		Header: http.Header{
			"Content-Type":        []string{"text/csv"},
			"Content-Disposition": []string{"attachment; filename=test.csv"},
		},
		StatusCode: http.StatusOK,
	}

	mockClient := mockApiClient{downloadResponse: download}

	handler := downloadCallback(mockClient)
	handler.ServeHTTP(w, r)
	resp := w.Result()
	respBody, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/csv", resp.Header.Get("Content-Type"))
	assert.Equal(t, "attachment; filename=test.csv", resp.Header.Get("Content-Disposition"))
	assert.Equal(t, fileContent, string(respBody))
}

func TestDownloadCallback_error(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "download/callback?uid=dGVzdC5jc3Y=", nil)

	download := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(bytes.NewReader([]byte{})),
	}

	mockClient := mockApiClient{downloadResponse: download}

	handler := downloadCallback(mockClient)
	handler.ServeHTTP(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
