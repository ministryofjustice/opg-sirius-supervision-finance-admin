package api

import (
	"github.com/opg-sirius-finance-admin/apierror"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
)

func (suite *IntegrationSuite) TestCheckDownload() {
	conn := suite.testDB.GetConn()

	req := httptest.NewRequest(http.MethodHead, "/download?uid=eyJLZXkiOiJ0ZXN0LmNzdiIsIlZlcnNpb25JZCI6InZwckF4c1l0TFZzYjVQOUhfcUhlTlVpVTlNQm5QTmN6In0=", nil)
	w := httptest.NewRecorder()

	mockS3 := MockFileStorage{}
	mockS3.exists = true

	server := NewServer(nil, conn.Conn, nil, &mockS3)
	err := server.checkDownload(w, req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *IntegrationSuite) TestCheckDownload_noMatch() {
	conn := suite.testDB.GetConn()

	req := httptest.NewRequest(http.MethodHead, "/download?uid=eyJLZXkiOiJ0ZXN0LmNzdiIsIlZlcnNpb25JZCI6InZwckF4c1l0TFZzYjVQOUhfcUhlTlVpVTlNQm5QTmN6In0=", nil)
	w := httptest.NewRecorder()

	mockS3 := MockFileStorage{}
	mockS3.exists = false

	server := NewServer(nil, conn.Conn, nil, &mockS3)
	err := server.checkDownload(w, req)

	assert.ErrorIs(suite.T(), err, apierror.NotFound{})
}
