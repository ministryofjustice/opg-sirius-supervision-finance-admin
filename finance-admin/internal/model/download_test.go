package model

import (
	"github.com/opg-sirius-finance-admin/shared"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewDownload(t *testing.T) {
	type args struct {
		reportType         string
		reportJournalType  string
		reportScheduleType string
		reportAccountType  string
		reportDebtType     string
		dateOfTransaction  string
		dateTo             string
		dateFrom           string
		email              string
	}

	dateOfTransaction, _ := time.Parse("02/01/2006", "11/05/2024")
	dateTo, _ := time.Parse("02/01/2006", "15/06/2025")
	dateFrom, _ := time.Parse("02/01/2006", "21/07/2022")

	tests := []struct {
		name string
		args args
		want Download
	}{
		{
			name: "Returns all fields",
			args: args{
				reportType:         "reportType",
				reportJournalType:  "reportJournalType",
				reportScheduleType: "reportScheduleType",
				reportAccountType:  "reportAccountType",
				reportDebtType:     "reportDebtType",
				dateOfTransaction:  "11/05/2024",
				dateTo:             "15/06/2025",
				dateFrom:           "21/07/2022",
				email:              "Something@example.com",
			},
			want: Download{
				ReportType:         "reportType",
				ReportJournalType:  "reportJournalType",
				ReportScheduleType: "reportScheduleType",
				ReportAccountType:  "reportAccountType",
				ReportDebtType:     "reportDebtType",
				DateOfTransaction:  &shared.Date{Time: dateOfTransaction},
				ToDateField:        &shared.Date{Time: dateTo},
				FromDateField:      &shared.Date{Time: dateFrom},
				Email:              "Something@example.com",
			},
		},
		{
			name: "Returns with missing optional fields",
			args: args{
				reportType:         "reportType",
				reportJournalType:  "reportJournalType",
				reportScheduleType: "reportScheduleType",
				reportAccountType:  "reportAccountType",
				reportDebtType:     "reportDebtType",
				dateOfTransaction:  "",
				dateTo:             "",
				dateFrom:           "",
				email:              "Something@example.com",
			},
			want: Download{
				ReportType:         "reportType",
				ReportJournalType:  "reportJournalType",
				ReportScheduleType: "reportScheduleType",
				ReportAccountType:  "reportAccountType",
				ReportDebtType:     "reportDebtType",
				DateOfTransaction:  nil,
				ToDateField:        nil,
				FromDateField:      nil,
				Email:              "Something@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDownload(
				tt.args.reportType,
				tt.args.reportJournalType,
				tt.args.reportScheduleType,
				tt.args.reportAccountType,
				tt.args.reportDebtType,
				tt.args.dateOfTransaction,
				tt.args.dateTo,
				tt.args.dateFrom,
				tt.args.email,
			)
			assert.Equal(t, tt.want, got)
		})
	}
}
