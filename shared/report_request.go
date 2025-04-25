package shared

import (
	"strconv"
)

type ReportRequest struct {
	ReportType             ReportsType             `json:"reportType"`
	JournalType            *JournalType            `json:"journalType,omitempty"`
	ScheduleType           *ScheduleType           `json:"scheduleType,omitempty"`
	AccountsReceivableType *AccountsReceivableType `json:"AccountsReceivableType,omitempty"`
	DebtType               *DebtType               `json:"debtType,omitempty"`
	TransactionDate        *Date                   `json:"transactionDate,omitempty"`
	ToDate                 *Date                   `json:"toDate,omitempty"`
	FromDate               *Date                   `json:"fromDate,omitempty"`
	Email                  string                  `json:"email"`
	PisNumber              int                     `json:"pisNumber"`
}

func NewReportRequest(reportType, journalType, scheduleType, accountsReceivableType, debtType, transactionDate, dateTo, dateFrom, email, pisNumber string) ReportRequest {
	download := ReportRequest{
		ReportType:             ParseReportsType(reportType),
		JournalType:            ParseJournalType(journalType),
		ScheduleType:           ParseScheduleType(scheduleType),
		AccountsReceivableType: ParseAccountsReceivableType(accountsReceivableType),
		DebtType:               ParseReportDebtType(debtType),
		Email:                  email,
	}

	if transactionDate != "" {
		raisedDateFormatted := NewDate(transactionDate)
		download.TransactionDate = &raisedDateFormatted
	}

	if dateTo != "" {
		startDateFormatted := NewDate(dateTo)
		download.ToDate = &startDateFormatted
	}

	if dateFrom != "" {
		endDateFormatted := NewDate(dateFrom)
		download.FromDate = &endDateFormatted
	}

	if pisNumber != "" {
		download.PisNumber, _ = strconv.Atoi(pisNumber)
	}

	return download
}
