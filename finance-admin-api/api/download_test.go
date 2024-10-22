package api

import (
	"bytes"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/opg-sirius-finance-admin/apierror"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_download(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/download/test.csv", nil)
	w := httptest.NewRecorder()
	req.SetPathValue("filename", "abc.csv")

	fileContent := "col1,col2,col3\n1,a,Z\n"

	mockAwsClient := MockAWSClient{}
	mockAwsClient.outgoingObject = &s3.GetObjectOutput{
		Body:        io.NopCloser(bytes.NewReader([]byte(fileContent))),
		ContentType: aws.String("text/csv"),
	}

	server := Server{&mockAwsClient}
	_ = server.download(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, fileContent, w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, res.Header.Get("Content-Type"), "text/csv")
	assert.Equal(t, res.Header.Get("Content-Disposition"), "attachment; filename=abc.csv")
}

func TestServer_download_noMatch(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/download/test.csv", nil)
	w := httptest.NewRecorder()
	req.SetPathValue("filename", "abc.csv")

	mockAwsClient := MockAWSClient{}
	mockAwsClient.err = &types.NoSuchKey{}

	server := Server{&mockAwsClient}
	err := server.download(w, req)

	expected := apierror.NotFoundError(&types.NoSuchKey{})
	assert.ErrorAs(t, err, &expected)
}