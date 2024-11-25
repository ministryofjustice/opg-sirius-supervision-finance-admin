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
)

func (suite *IntegrationSuite) TestServer_download() {
	conn := suite.testDB.GetConn()

	req := httptest.NewRequest(http.MethodGet, "/download?uid=eyJLZXkiOiJ0ZXN0LmNzdiIsIlZlcnNpb25JZCI6InZwckF4c1l0TFZzYjVQOUhfcUhlTlVpVTlNQm5QTmN6In0=", nil)
	w := httptest.NewRecorder()

	fileContent := "col1,col2,col3\n1,a,Z\n"

	mockS3 := MockFileStorage{}
	mockS3.outgoingObject = &s3.GetObjectOutput{
		Body:        io.NopCloser(bytes.NewBufferString(fileContent)),
		ContentType: aws.String("text/csv"),
	}

	server := NewServer(nil, conn.Conn, nil, &mockS3)
	_ = server.download(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(suite.T(), fileContent, w.Body.String())
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), res.Header.Get("Content-Type"), "text/csv")
	assert.Equal(suite.T(), res.Header.Get("Content-Disposition"), "attachment; filename=test.csv")
}

func (suite *IntegrationSuite) TestServer_download_noMatch() {
	conn := suite.testDB.GetConn()

	req := httptest.NewRequest(http.MethodGet, "/download?uid=eyJLZXkiOiJ0ZXN0LmNzdiIsIlZlcnNpb25JZCI6InZwckF4c1l0TFZzYjVQOUhfcUhlTlVpVTlNQm5QTmN6In0=", nil)
	w := httptest.NewRecorder()

	mockS3 := MockFileStorage{}
	mockS3.err = &types.NoSuchKey{}
	server := NewServer(nil, conn.Conn, nil, &mockS3)

	err := server.download(w, req)

	expected := apierror.NotFoundError(&types.NoSuchKey{})
	assert.ErrorAs(suite.T(), err, &expected)
}
