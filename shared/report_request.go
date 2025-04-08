package shared

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
	PisNumber              *int                    `json:"pisNumber,omitempty"`
}

func NewReportRequest(reportType, journalType, scheduleType, accountsReceivableType, debtType, transactionDate, dateTo, dateFrom, email string, pisNumber *int) ReportRequest {
	download := ReportRequest{
		ReportType:             ParseReportsType(reportType),
		JournalType:            ParseJournalType(journalType),
		ScheduleType:           ParseScheduleType(scheduleType),
		AccountsReceivableType: ParseAccountsReceivableType(accountsReceivableType),
		DebtType:               ParseReportDebtType(debtType),
		Email:                  email,
		PisNumber:              pisNumber,
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
