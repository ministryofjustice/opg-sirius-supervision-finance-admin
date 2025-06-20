package api

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestDownload(t *testing.T) {
	mockClient := &MockClient{}
	mockJwtClient := &mockJWTClient{}
	client := NewClient(mockClient, mockJwtClient, EnvVars{"http://localhost:3000", ""})
	fileContent := []byte("file content")

	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       io.NopCloser(bytes.NewReader(fileContent)),
		}, nil
	}

	resp, err := client.Download(testContext(), "dGVzdC5jc3Y=")
	assert.NoError(t, err)

	actual, _ := io.ReadAll(io.NopCloser(resp.Body))
	assert.Equal(t, fileContent, actual)
}
