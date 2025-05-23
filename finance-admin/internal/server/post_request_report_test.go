package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/api"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRequestReportHandlerSuccess(t *testing.T) {
	form := url.Values{
		"reportType":             {"AccountsReceivable"},
		"reportJournalType":      {""},
		"reportScheduleType":     {""},
		"accountsReceivableType": {"BadDebtWriteOff"},
		"reportDebtType":         {""},
		"dateOfTransaction":      {"11/05/2024"},
		"dateFrom":               {"01/04/2024"},
		"dateTo":                 {"31/03/2025"},
		"email":                  {"SomeSortOfEmail@example.com"},
	}

	client := mockApiClient{}
	ro := &mockRoute{client: client}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/download", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	appVars := AppVars{
		Path: "/download",
	}

	appVars.EnvironmentVars.Prefix = "prefix"

	sut := RequestReportHandler{ro}

	err := sut.render(appVars, w, r)

	assert.Nil(t, err)
}

func TestRequestReportHandlerValidationErrors(t *testing.T) {
	assert := assert.New(t)
	client := &mockApiClient{}
	ro := &mockRoute{client: client}

	validationErrors := model.ValidationErrors{
		"ReportType": {
			"ReportType": "Please select a report type",
		},
	}

	client.error = model.ValidationError{
		Errors: validationErrors,
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/download", nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	appVars := AppVars{
		Path: "/add",
	}

	sut := RequestReportHandler{ro}
	err := sut.render(appVars, w, r)
	assert.Nil(err)
	assert.Equal("422 Unprocessable Entity", w.Result().Status)
}

func TestRequestReportHandlerStatusError(t *testing.T) {
	assert := assert.New(t)
	client := &mockApiClient{}
	ro := &mockRoute{client: client}

	client.error = api.StatusError{
		Code:   http.StatusInternalServerError,
		URL:    "/downloads",
		Method: http.MethodGet,
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/download", nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	appVars := AppVars{
		Path: "/add",
	}

	sut := RequestReportHandler{ro}
	err := sut.render(appVars, w, r)
	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, w.Result().StatusCode)
}
