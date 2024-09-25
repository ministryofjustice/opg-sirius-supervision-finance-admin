package model

import "github.com/opg-sirius-finance-admin/shared"

type Download struct {
	ReportType         string       `json:"reportType"`
	ReportJournalType  string       `json:"reportJournalType"`
	ReportScheduleType string       `json:"reportScheduleType"`
	ReportAccountType  string       `json:"reportAccountType"`
	ReportDebtType     string       `json:"reportDebtType"`
	DateOfTransaction  *shared.Date `json:"dateOfTransaction,omitempty"`
	ToDateField        *shared.Date `json:"toDateField,omitempty"`
	FromDateField      *shared.Date `json:"fromDateField,omitempty"`
	Email              string       `json:"email"`
}

func NewDownload(reportType, reportJournalType, reportScheduleType, reportAccountType, reportDebtType, dateOfTransaction, dateTo, dateFrom, email string) Download {
	download := Download{
		ReportType:         reportType,
		ReportJournalType:  reportJournalType,
		ReportScheduleType: reportScheduleType,
		ReportAccountType:  reportAccountType,
		ReportDebtType:     reportDebtType,
		Email:              email,
	}

	if dateOfTransaction != "" {
		raisedDateFormatted := shared.NewDate(dateOfTransaction)
		download.DateOfTransaction = &raisedDateFormatted
	}

	if dateTo != "" {
		startDateFormatted := shared.NewDate(dateTo)
		download.ToDateField = &startDateFormatted
	}

	if dateFrom != "" {
		endDateFormatted := shared.NewDate(dateFrom)
		download.FromDateField = &endDateFormatted
	}

	return download
}
