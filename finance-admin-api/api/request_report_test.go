package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/opg-sirius-finance-admin/apierror"
	"github.com/opg-sirius-finance-admin/shared"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRequestReport(t *testing.T) {
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
	mockDb := MockDb{}

	server := Server{&mockHttpClient, &mockDb, &mockDispatch, &mockFileStorage}
	_ = server.requestReport(w, req)

	res := w.Result()
	defer res.Body.Close()

	expected := ""

	assert.Equal(t, expected, w.Body.String())
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRequestReportNoEmail(t *testing.T) {
	var b bytes.Buffer

	downloadForm := shared.Download{
		ReportType:        "AccountsReceivable",
		ReportAccountType: "AgedDebt",
		Email:             "",
	}

	_ = json.NewEncoder(&b).Encode(downloadForm)
	req := httptest.NewRequest(http.MethodPost, "/downloads", &b)
	w := httptest.NewRecorder()

	mockHttpClient := MockHttpClient{}
	mockDispatch := MockDispatch{}
	mockFileStorage := MockFileStorage{}

	server := Server{&mockHttpClient, nil, &mockDispatch, &mockFileStorage}
	err := server.requestReport(w, req)

	res := w.Result()
	defer res.Body.Close()

	expected := apierror.ValidationError{Errors: apierror.ValidationErrors{
		"Email": {
			"required": "This field Email needs to be looked at required",
		},
	},
	}

	assert.Equal(t, expected, err)
}
func TestGenerateAndUploadReport(t *testing.T) {
	mockHttpClient := MockHttpClient{}
	mockDispatch := MockDispatch{}
	mockFileStorage := MockFileStorage{}
	mockDb := MockDb{}

	mockFileStorage.versionId = "123"

	server := Server{&mockHttpClient, &mockDb, &mockDispatch, &mockFileStorage}

	ctx := context.Background()
	timeNow, _ := time.Parse("2006-01-02", "2024-01-01")
	toDate := shared.NewDate("2024-01-01")
	fromDate := shared.NewDate("2024-10-01")
	download := shared.Download{
		ReportAccountType: "AgedDebt",
		ToDateField:       &toDate,
		FromDateField:     &fromDate,
	}

	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}, nil
	}

	err := server.generateAndUploadReport(ctx, download, timeNow)

	assert.Equal(t, nil, err)
}
