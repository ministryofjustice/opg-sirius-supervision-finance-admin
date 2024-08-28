package model

type Upload struct {
	ReportUploadType string `json:"reportUploadType"`
	UploadDate       *Date  `json:"uploadDate"`
	Email            string `json:"email"`
	File             []byte `json:"file"`
}
