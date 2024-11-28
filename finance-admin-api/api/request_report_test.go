package api

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/apierror"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/shared"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestRequestReport(t *testing.T) {
	var b bytes.Buffer

	downloadForm := &shared.ReportRequest{
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

	downloadForm := shared.ReportRequest{
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
	download := shared.ReportRequest{
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

func TestCreateCsv(t *testing.T) {
	want, _ := os.Create("test.csv")
	defer want.Close()

	writer := csv.NewWriter(want)
	_ = writer.Write([]string{"test", "hehe"})
	_ = writer.Write([]string{"123 Real Street", "Bingopolis"})
	writer.Flush()

	items := [][]string{{"test", "hehe"}, {"123 Real Street", "Bingopolis"}}
	_, err := createCsv("test2.csv", items)

	wantBytes, _ := os.ReadFile("test.csv")
	gotBytes, _ := os.ReadFile("test2.csv")

	assert.Nil(t, err)
	assert.Equal(t, string(wantBytes), string(gotBytes))
}

func TestCreateCsvNoItems(t *testing.T) {
	items := [][]string{}
	_, err := createCsv("test.csv", items)
	gotBytes, _ := os.ReadFile("test.csv")

	assert.Nil(t, err)
	assert.Equal(t, "", string(gotBytes))
}

func TestCreateDownloadNotifyPayload(t *testing.T) {
	emailAddress := "test@email.com"
	reportName := "test report"
	downloadRequest := shared.DownloadRequest{
		Key:       "test.csv",
		VersionId: "1",
	}
	uid, _ := downloadRequest.Encode()
	requestedDate, _ := time.Parse("2006-01-02 15:04:05", "2024-01-01 13:37:00")
	os.Setenv("SIRIUS_PUBLIC_URL", "www.sirius.com")
	os.Setenv("PREFIX", "/finance")

	want := NotifyPayload{
		EmailAddress: emailAddress,
		TemplateId:   reportRequestedTemplateId,
		Personalisation: reportRequestedNotifyPersonalisation{
			FileLink:          fmt.Sprintf("www.sirius.com/finance/download?uid=%s", uid),
			ReportName:        reportName,
			RequestedDate:     "2024-01-01",
			RequestedDateTime: "2024-01-01 13:37:00",
		},
	}

	payload, err := createDownloadNotifyPayload(emailAddress, downloadRequest.Key, &downloadRequest.VersionId, requestedDate, reportName)

	assert.Equal(t, want, payload)
	assert.Nil(t, err)
}
