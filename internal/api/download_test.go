package api

import (
	"bytes"
	"encoding/json"
	"github.com/opg-sirius-finance-admin/internal/model"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSubmitDownload(t *testing.T) {
	mockClient := &MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000", "")

	data := `{
		"reportType":         "AccountsReceivable",
		"reportJournalType":  "",
		"reportScheduleType": "",
		"reportAccountType":  "BadDebtWriteOffReport",
		"reportDebtType":     "",
		"dateOfTransaction":  "11/05/2024",
		"dateFrom":           "01/04/2024",
		"dateTo":             "31/03/2025",
		"email":              "SomeSortOfEmail@example.com",
	}
	`

	r := io.NopCloser(bytes.NewReader([]byte(data)))

	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 201,
			Body:       r,
		}, nil
	}

	err := client.Download(getContext(nil), "AccountsReceivable", "", "", "BadDebtWriteOffReport", "", "11/05/2024", "01/04/2024", "31/03/2025", "SomeSortOfEmail@example.com")
	assert.Equal(t, nil, err)
}

func TestSubmitDownloadUnauthorised(t *testing.T) {
	mockClient := &MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000", "")

	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(bytes.NewReader([]byte{})), // Empty body
		}, nil
	}

	err := client.Download(getContext(nil), "AccountsReceivable", "", "", "BadDebtWriteOffReport", "", "11/05/2024", "01/04/2024", "31/03/2025", "SomeSortOfEmail@example.com")

	assert.Equal(t, ErrUnauthorized.Error(), err.Error())
}

func TestSubmitDownloadReturnsBadRequestError(t *testing.T) {
	mockClient := &MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000", "")

	json := `
		{"reasons":["StartDate","EndDate"]}
	`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 400,
			Body:       r,
		}, nil
	}

	err := client.Download(getContext(nil), "AccountsReceivable", "", "", "BadDebtWriteOffReport", "", "11/05/2024", "01/04/2024", "31/03/2025", "SomeSortOfEmail@example.com")

	expectedError := model.ValidationError{Message: "", Errors: model.ValidationErrors{"EndDate": map[string]string{"EndDate": "EndDate"}, "StartDate": map[string]string{"StartDate": "StartDate"}}}
	assert.Equal(t, expectedError, err)
}

func TestSubmitDownloadReturnsValidationError(t *testing.T) {
	validationErrors := model.ValidationError{
		Message: "Validation failed",
		Errors: map[string]map[string]string{
			"ReportType": {
				"required": "Please select a report type",
			},
		},
	}
	responseBody, _ := json.Marshal(validationErrors)
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write(responseBody)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL, svr.URL)

	err := client.Download(getContext(nil), "", "", "", "", "", "", "", "", "")
	expectedError := model.ValidationError{Message: "", Errors: model.ValidationErrors{"ReportType": map[string]string{"required": "Please select a report type"}}}
	assert.Equal(t, expectedError, err.(model.ValidationError))
}

func TestNewDownload(t *testing.T) {
	type args struct {
		dateOfTransaction string
		dateTo            string
		dateFrom          string
	}
	tests := []struct {
		name  string
		args  args
		want  *model.Date
		want1 *model.Date
		want2 *model.Date
	}{
		{
			name: "No dates passed in no dates are returned",
			args: args{
				dateOfTransaction: "",
				dateTo:            "",
				dateFrom:          "",
			},
			want:  nil,
			want1: nil,
			want2: nil,
		},
		{
			name: "All dates passed in all dates are returned",
			args: args{
				dateOfTransaction: "01/01/2021",
				dateTo:            "02/02/2022",
				dateFrom:          "03/03/2023",
			},
			want:  &model.Date{Time: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)},
			want1: &model.Date{Time: time.Date(2022, time.February, 2, 0, 0, 0, 0, time.UTC)},
			want2: &model.Date{Time: time.Date(2023, time.March, 3, 0, 0, 0, 0, time.UTC)},
		},
		{
			name: "Only one date passed in one date is returned",
			args: args{
				dateOfTransaction: "01/01/2021",
				dateTo:            "",
				dateFrom:          "",
			},
			want:  &model.Date{Time: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)},
			want1: nil,
			want2: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := NewDownload(tt.args.dateOfTransaction, tt.args.dateTo, tt.args.dateFrom)
			assert.Equalf(t, tt.want, got, "NewDownload(%v, %v, %v)", tt.args.dateOfTransaction, tt.args.dateTo, tt.args.dateFrom)
			assert.Equalf(t, tt.want1, got1, "NewDownload(%v, %v, %v)", tt.args.dateOfTransaction, tt.args.dateTo, tt.args.dateFrom)
			assert.Equalf(t, tt.want2, got2, "NewDownload(%v, %v, %v)", tt.args.dateOfTransaction, tt.args.dateTo, tt.args.dateFrom)
		})
	}
}
