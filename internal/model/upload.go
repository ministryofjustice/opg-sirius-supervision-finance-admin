package model

import "os"

type Upload struct {
	ReportUploadType string   `json:"reportUploadType"`
	UploadDate       *Date    `json:"uploadDate"`
	Email            string   `json:"email"`
	File             *os.File `json:"file"`
}
