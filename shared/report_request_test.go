package shared

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewReportRequest(t *testing.T) {
	type args struct {
		reportType                   string
		reportJournalType            string
		reportScheduleType           string
		ReportAccountsReceivableType string
		reportDebtType               string
		dateOfTransaction            string
		dateTo                       string
		dateFrom                     string
		pisNumber                    string
		email                        string
	}

	dateOfTransaction, _ := time.Parse("02/01/2006", "11/05/2024")
	dateTo, _ := time.Parse("02/01/2006", "15/06/2025")
	dateFrom, _ := time.Parse("02/01/2006", "21/07/2022")

	tests := []struct {
		name string
		args args
		want ReportRequest
	}{
		{
			name: "Returns all fields",
			args: args{
				reportType:                   ReportsTypeSchedule.String(),
				reportJournalType:            JournalTypeReceiptTransactions.String(),
				reportScheduleType:           ScheduleTypeOnlineCardPayments.String(),
				ReportAccountsReceivableType: AccountsReceivableTypeAgedDebt.String(),
				reportDebtType:               DebtTypeFeeChase.String(),
				dateOfTransaction:            "11/05/2024",
				dateTo:                       "15/06/2025",
				dateFrom:                     "21/07/2022",
				pisNumber:                    "123456",
				email:                        "Something@example.com",
			},
			want: ReportRequest{
				ReportType:             ReportsTypeSchedule,
				JournalType:            toPtr(JournalTypeReceiptTransactions),
				ScheduleType:           toPtr(ScheduleTypeOnlineCardPayments),
				AccountsReceivableType: toPtr(AccountsReceivableTypeAgedDebt),
				DebtType:               toPtr(DebtTypeFeeChase),
				TransactionDate:        &Date{Time: dateOfTransaction},
				ToDate:                 &Date{Time: dateTo},
				FromDate:               &Date{Time: dateFrom},
				PisNumber:              toPtr(123456),
				Email:                  "Something@example.com",
			},
		},
		{
			name: "Returns with missing optional fields",
			args: args{
				reportType:                   ReportsTypeSchedule.String(),
				reportJournalType:            "",
				reportScheduleType:           "",
				ReportAccountsReceivableType: "",
				reportDebtType:               "",
				dateOfTransaction:            "",
				dateTo:                       "",
				dateFrom:                     "",
				pisNumber:                    "",
				email:                        "Something@example.com",
			},
			want: ReportRequest{
				ReportType:             ReportsTypeSchedule,
				JournalType:            nil,
				ScheduleType:           nil,
				AccountsReceivableType: nil,
				DebtType:               nil,
				TransactionDate:        nil,
				ToDate:                 nil,
				FromDate:               nil,
				PisNumber:              nil,
				Email:                  "Something@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewReportRequest(
				tt.args.reportType,
				tt.args.reportJournalType,
				tt.args.reportScheduleType,
				tt.args.ReportAccountsReceivableType,
				tt.args.reportDebtType,
				tt.args.dateOfTransaction,
				tt.args.dateTo,
				tt.args.dateFrom,
				tt.args.pisNumber,
				tt.args.email,
			)
			assert.Equal(t, tt.want, got)
		})
	}
}

func toPtr[T any](val T) *T {
	return &val
}
