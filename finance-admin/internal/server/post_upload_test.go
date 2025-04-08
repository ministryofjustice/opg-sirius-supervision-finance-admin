package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestUploadFormHandlerNoFileUploaded(t *testing.T) {
	form := url.Values{
		"reportUploadType": {"DEBT_CHASE"},
		"uploadDate":       {"2024-01-24"},
		"email":            {"SomeSortOfEmail@example.com"},
	}

	client := mockApiClient{}
	ro := &mockRoute{client: client}

	sut := UploadFormHandler{ro}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/uploads", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	appVars := AppVars{
		Path: "/uploads",
	}
	appVars.EnvironmentVars.Prefix = "prefix"

	err := sut.render(appVars, w, r)
	_, _ = w.Write([]byte("No file uploaded"))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "No file uploaded")
	assert.Nil(t, err)
}

func TestUploadFormHandlerSuccess(t *testing.T) {
	form := url.Values{
		"reportUploadType": {"DEBT_CHASE"},
		"uploadDate":       {"2024-01-24"},
		"email":            {"SomeSortOfEmail@example.com"},
	}

	client := mockApiClient{}
	ro := &mockRoute{client: client}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/uploads", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	appVars := AppVars{
		Path: "/uploads",
	}

	appVars.EnvironmentVars.Prefix = "prefix"
	sut := UploadFormHandler{ro}
	err := sut.render(appVars, w, r)
	assert.Nil(t, err)
}

func TestUploadFormHandlerValidationErrors(t *testing.T) {
	assert := assert.New(t)
	client := &mockApiClient{}
	ro := &mockRoute{client: client}

	validationErrors := model.ValidationErrors{
		"ReportUploadType": {
			"ReportUploadType": "Please select a report type",
		},
	}

	client.error = model.ValidationError{
		Errors: validationErrors,
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/uploads", nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	appVars := AppVars{
		Path: "/uploads",
	}

	sut := UploadFormHandler{ro}
	err := sut.render(appVars, w, r)
	assert.Nil(err)
	assert.Equal("400 Bad Request", w.Result().Status)
}

func TestUploadDateNotInTheFutureValidationErrors(t *testing.T) {
	assert := assert.New(t)
	client := &mockApiClient{}
	ro := &mockRoute{client: client}

	validationErrors := model.ValidationErrors{
		"UploadDate": {
			"date-in-the-future": "The report date must be today or in the past",
		},
	}

	client.error = model.ValidationError{
		Errors: validationErrors,
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/uploads", nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	appVars := AppVars{
		Path: "/uploads",
	}

	sut := UploadFormHandler{ro}
	err := sut.render(appVars, w, r)
	assert.Nil(err)
	assert.Equal("400 Bad Request", w.Result().Status)
}
