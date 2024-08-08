package model

type Download struct {
	ReportType         string `json:"reportType"`
	ReportJournalType  string `json:"reportJournalType"`
	ReportScheduleType string `json:"reportScheduleType"`
	ReportAccountType  string `json:"reportAccountType"`
	ReportDebtType     string `json:"reportDebtType"`
	DateField          *Date  `json:"dateField,omitempty"`
	ToDateField        *Date  `json:"toDateField,omitempty"`
	FromDateField      *Date  `json:"fromDateField,omitempty"`
	Email              string `json:"email"`
}
