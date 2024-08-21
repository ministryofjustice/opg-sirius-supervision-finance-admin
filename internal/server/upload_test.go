package server

import (
	"bytes"
	"github.com/opg-sirius-finance-admin/internal/model"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestUploadNoFileUploaded(t *testing.T) {
	form := url.Values{
		"reportUploadType": {"DebtChase"},
		"uploadDate":       {"2024-01-24"},
		"email":            {"SomeSortOfEmail@example.com"},
	}

	client := mockApiClient{}
	ro := &mockRoute{client: client}

	sut := UploadHandler{ro}

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

func TestUploadFailedToReadCSVHeaders(t *testing.T) {
	form := url.Values{
		"reportUploadType": {"DebtChase"},
		"uploadDate":       {"2024-01-24"},
		"email":            {"SomeSortOfEmail@example.com"},
	}

	mockCSVData := "Header1,Header2\nValue1,Value2\n"

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	filePart, err := writer.CreateFormFile("fileUpload", "mock.csv")
	defer func() {
		os.Remove("mock.csv")
	}()

	assert.NoError(t, err)
	_, err = filePart.Write([]byte(mockCSVData))
	assert.NoError(t, err)

	for key, values := range form {
		for _, value := range values {
			_ = writer.WriteField(key, value)
		}
	}

	err = writer.Close()
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/uploads", &buf)
	r.Header.Set("Content-Type", writer.FormDataContentType())

	client := mockApiClient{}
	ro := &mockRoute{client: client}
	sut := UploadHandler{ro}

	_ = sut.render(AppVars{}, w, r)
	_, _ = w.Write([]byte("Failed to read CSV headers"))
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to read CSV headers")
}

func TestUploadSuccess(t *testing.T) {
	form := url.Values{
		"reportUploadType": {"DebtChase"},
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
	sut := UploadHandler{ro}
	err := sut.render(appVars, w, r)
	assert.Nil(t, err)
}

func TestUploadValidationErrors(t *testing.T) {
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

	sut := UploadHandler{ro}
	err := sut.render(appVars, w, r)
	assert.Nil(err)
	assert.Equal("400 Bad Request", w.Result().Status)
}

func Test_reportHeadersByType(t *testing.T) {
	tests := []struct {
		name       string
		reportType string
		want       []string
	}{
		{
			name:       "Deputy schedule report get the correct header",
			reportType: "DeputySchedule",
			want:       []string{"Deputy number", "Deputy name", "Case number", "Client forename", "Client surname", "Do not invoice", "Total outstanding"},
		},
		{
			name:       "Debt chase report get the correct header",
			reportType: "DebtChase",
			want:       []string{"Client_no", "Deputy_name", "Total_debt"},
		},
		{
			name:       "Payments OPG BACS report get the correct header",
			reportType: "PaymentsOPGBACS",
			want:       []string{"Line", "Type", "Code", "Number", "Transaction", "Value Date", "Amount", "Amount Reconciled", "Charges", "Status", "Desc Flex", "Consolidated line"},
		},
		{
			name:       "No match will return unknown",
			reportType: "",
			want:       []string{"Unknown report type"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, reportHeadersByType(tt.reportType), "reportHeadersByType(%v)", tt.reportType)
		})
	}
}

func TestCleanString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Removes white space from start and end", "  Hello, World!  ", "Hello, World!"},
		{"Removes new lines and tabs", "\n\tHello, World!\n\t", "Hello, World!"},
		{"Removes nil character", "Hello,\x00World!", "Hello,World!"},
		{"Nothing is removed", "", ""}, // empty string should return empty string
		{"Remove only whitespace and control characters", "  \t\n  \x0B\x0C   ", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := cleanString(tt.input)
			if got != tt.expected {
				t.Errorf("cleanString(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}
