package shared

type ReportRequest struct {
	ReportType             ReportsType                  `json:"reportType"`
	JournalType            ReportJournalType            `json:"journalType"`
	ScheduleType           ReportScheduleType           `json:"scheduleType"`
	AccountsReceivableType ReportAccountsReceivableType `json:"AccountsReceivableType"`
	DebtType               ReportDebtType               `json:"debtType"`
	TransactionDate        *Date                        `json:"transactionDate,omitempty"`
	ToDate                 *Date                        `json:"toDate,omitempty"`
	FromDate               *Date                        `json:"fromDate,omitempty"`
	Email                  string                       `json:"email"`
}

func NewReportRequest(reportType, journalType, scheduleType, accountsReceivableType, debtType, transactionDate, dateTo, dateFrom, email string) ReportRequest {
	download := ReportRequest{
		ReportType:             ParseReportsType(reportType),
		JournalType:            ParseReportJournalType(journalType),
		ScheduleType:           ParseReportScheduleType(scheduleType),
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

	return download
}
