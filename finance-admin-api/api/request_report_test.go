package api

import (
	"bytes"
	"encoding/json"
	"github.com/opg-sirius-finance-admin/shared"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
)

func (suite *IntegrationSuite) TestRequestReport() {
	conn := suite.testDB.GetConn()

	var b bytes.Buffer

	downloadForm := &shared.Download{
		ReportType:        "AccountsReceivable",
		ReportAccountType: "AgedDebt",
		Email:             "joseph@test.com",
	}

	_ = json.NewEncoder(&b).Encode(downloadForm)
	req := httptest.NewRequest(http.MethodPost, "/downloads", &b)
	w := httptest.NewRecorder()

	mockHttpClient := MockHttpClient{}
	mockDispatch := MockDispatch{}
	mockFileStorage := MockFileStorage{}

	server := Server{&mockHttpClient, conn.Conn, &mockDispatch, &mockFileStorage}
	_ = server.requestReport(w, req)

	res := w.Result()
	defer res.Body.Close()

	expected := ""

	assert.Equal(suite.T(), expected, w.Body.String())
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
}
