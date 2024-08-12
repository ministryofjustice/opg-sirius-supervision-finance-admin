package server

import (
	"github.com/opg-sirius-finance-admin/internal/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestDownloadSuccess(t *testing.T) {
	form := url.Values{
		"reportType":         {"AccountsReceivable"},
		"reportJournalType":  {""},
		"reportScheduleType": {""},
		"reportAccountType":  {"BadDebtWriteOffReport"},
		"reportDebtType":     {""},
		"dateField":          {"11/05/2024"},
		"dateFromField":      {"01/04/2024"},
		"dateToField":        {"31/03/2025"},
		"emailField":         {"SomeSortOfEmail@example.com"},
	}

	client := mockApiClient{}
	ro := &mockRoute{client: client}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/download", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.SetPathValue("clientId", "1")

	appVars := AppVars{
		Path: "/download",
	}

	appVars.EnvironmentVars.Prefix = "prefix"

	sut := DownloadHandler{ro}

	err := sut.render(appVars, w, r)

	assert.Nil(t, err)
}

func TestDownloadValidationErrors(t *testing.T) {
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
	r.SetPathValue("clientId", "1")

	appVars := AppVars{
		Path: "/add",
	}

	sut := DownloadHandler{ro}
	err := sut.render(appVars, w, r)
	assert.Nil(err)
	assert.Equal("422 Unprocessable Entity", w.Result().Status)
}
