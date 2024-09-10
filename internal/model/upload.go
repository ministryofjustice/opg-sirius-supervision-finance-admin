package model

import "io"

type Upload struct {
	ReportUploadType string `json:"reportUploadType"`
	UploadDate       *Date  `json:"uploadDate"`
	Email            string `json:"email"`
	File             []byte `json:"file"`
}

func NewUpload(reportUploadType string, uploadDate string, email string, file io.Reader) (Upload, error) {
	fileTransformed, err := io.ReadAll(file)
	if err != nil {
		return Upload{}, err
	}

	upload := Upload{
		ReportUploadType: reportUploadType,
		Email:            email,
		File:             fileTransformed,
	}

	if uploadDate != "" {
		uploadDateFormatted := NewDate(uploadDate)
		upload.UploadDate = &uploadDateFormatted
	}

	return upload, nil
}
