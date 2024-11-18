package shared

type Download struct {
	ReportType        string `json:"reportType"`
	ReportAccountType string `json:"reportAccountType"`
	DateOfTransaction string `json:"dateOfTransaction"`
}
