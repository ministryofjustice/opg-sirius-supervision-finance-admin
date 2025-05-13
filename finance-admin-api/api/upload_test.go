package api

import (
	"bytes"
	"encoding/json"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/apierror"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/shared"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_upload(t *testing.T) {
	var b bytes.Buffer

	uploadForm := &shared.Upload{
		ReportUploadType: shared.ReportTypeUploadDebtChase,
		Email:            "joseph@test.com",
		Filename:         "file.txt",
		File:             []byte("client_no,deputy_name,Total_debt"),
	}

	_ = json.NewEncoder(&b).Encode(uploadForm)
	req := httptest.NewRequest(http.MethodPost, "/uploads", &b)
	w := httptest.NewRecorder()

	mockHttpClient := MockHttpClient{}
	mockDispatch := MockDispatch{}
	mockFileStorage := MockFileStorage{}

	server := Server{&mockHttpClient, &mockDispatch, &mockFileStorage}
	_ = server.upload(w, req)

	res := w.Result()
	defer res.Body.Close()

	expected := ""

	assert.Equal(t, expected, w.Body.String())
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestUploadIncorrectCSVHeaders(t *testing.T) {
	var b bytes.Buffer

	uploadForm := &shared.Upload{
		ReportUploadType: shared.ReportTypeUploadDebtChase,
		Email:            "joseph@test.com",
		Filename:         "file.txt",
		File:             []byte("blarg"),
	}

	_ = json.NewEncoder(&b).Encode(uploadForm)
	req := httptest.NewRequest(http.MethodPost, "/uploads", &b)
	w := httptest.NewRecorder()

	mockHttpClient := MockHttpClient{}
	mockDispatch := MockDispatch{}

	server := Server{&mockHttpClient, &mockDispatch, nil}
	err := server.upload(w, req)

	expected := apierror.ValidationError{Errors: apierror.ValidationErrors{
		"FileUpload": {
			"incorrect-headers": "CSV headers do not match for the report trying to be uploaded",
		},
	}}
	assert.Equal(t, expected, err)
}

func TestUploadFailedToReadCSVHeaders(t *testing.T) {
	var b bytes.Buffer

	uploadForm := &shared.Upload{
		ReportUploadType: shared.ReportTypeUploadDebtChase,
		Email:            "joseph@test.com",
		Filename:         "file.txt",
		File:             []byte(""),
	}

	_ = json.NewEncoder(&b).Encode(uploadForm)
	req := httptest.NewRequest(http.MethodPost, "/uploads", &b)
	w := httptest.NewRecorder()

	mockHttpClient := MockHttpClient{}
	mockDispatch := MockDispatch{}

	server := Server{&mockHttpClient, &mockDispatch, nil}
	err := server.upload(w, req)

	expected := apierror.ValidationError{Errors: apierror.ValidationErrors{
		"FileUpload": {
			"read-failed": "Failed to read CSV headers",
		},
	}}
	assert.Equal(t, expected, err)
}

func TestUploadSkipHeaderValidation(t *testing.T) {
	var b bytes.Buffer

	uploadForm := &shared.Upload{
		ReportUploadType: shared.ReportTypeUploadDirectDebitsCollections,
		Email:            "joseph@test.com",
		Filename:         "file.txt",
		File:             []byte(""),
	}

	_ = json.NewEncoder(&b).Encode(uploadForm)
	req := httptest.NewRequest(http.MethodPost, "/uploads", &b)
	w := httptest.NewRecorder()

	mockHttpClient := MockHttpClient{}
	mockDispatch := MockDispatch{}
	mockFileStorage := MockFileStorage{}

	server := Server{&mockHttpClient, &mockDispatch, &mockFileStorage}
	_ = server.upload(w, req)

	res := w.Result()
	defer res.Body.Close()

	expected := ""

	assert.Equal(t, expected, w.Body.String())
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestValidateCSVHeaders(t *testing.T) {
	tests := []struct {
		name                string
		file                []byte
		reportType          shared.ReportUploadType
		useStrictComparison bool
		wantErr             bool
	}{
		{
			"Failed to read empty file",
			[]byte{},
			shared.ReportTypeUploadPaymentsSupervisionCheque,
			true,
			true,
		},
		{
			"Successfully validates only first row",
			[]byte("Line, Type, Code, Number, Transaction Date, Value Date, Amount, Amount Reconciled, Charges, Status, Desc Flex\nTest"),
			shared.ReportTypeUploadPaymentsOPGBACS,
			true,
			false,
		},
		{
			"Too many headers causes an error",
			[]byte("Line, Type, Code, Number, Transaction Date, Value Date, Amount, Amount Reconciled, Charges, Status, Desc Flex, Hehe\nTest"),
			shared.ReportTypeUploadPaymentsOPGBACS,
			true,
			true,
		},
		{
			"Un-strict comparison allows extra characters before and after headers",
			[]byte("1. Line (Defined by Sirius), 2. Type (£)\nTest"),
			shared.ReportTypeUploadPaymentsOPGBACS,
			false,
			false,
		},
		{
			"Ignores blank fields",
			[]byte("Line, Type,,,\nTest"),
			shared.ReportTypeUploadPaymentsOPGBACS,
			true,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCSVHeaders(tt.file, tt.reportType, tt.useStrictComparison)
			if tt.wantErr && err == nil {
				t.Errorf("wanted error got no error")
			} else if !tt.wantErr && err != nil {
				t.Errorf("wanted no error got error")
			}
		})
	}
}

func TestCleanString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Removes white space from start and end", "  Hello, World!  ", "hello, world!"},
		{"Removes new lines and tabs", "\n\tHello, World!\n\t", "hello, world!"},
		{"Removes nil character", "Hello,\x00World!", "hello,world!"},
		{"Nothing is removed", "", ""}, // empty string should return empty string
		{"Remove only whitespace and control characters", "  \t\n  \x0B\x0C   ", ""},
		{"Double space is replaced with single space", "Hello,  World!", "hello, world!"},
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
