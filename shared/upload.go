package shared

import "io"

type Upload struct {
	ReportUploadType ReportUploadType `json:"reportUploadType"`
	PisNumber        int              `json:"pisNumber"`
	UploadDate       Date             `json:"uploadDate"`
	Email            string           `json:"email"`
	Filename         string           `json:"filename"`
	File             []byte           `json:"file"`
}

func NewUpload(reportUploadType ReportUploadType, pisNumber int, uploadDate string, email string, file io.Reader, filename string) (Upload, error) {
	fileTransformed, err := io.ReadAll(file)
	if err != nil {
		return Upload{}, err
	}

	upload := Upload{
		ReportUploadType: reportUploadType,
		PisNumber:        pisNumber,
		Email:            email,
		File:             fileTransformed,
		Filename:         filename,
	}

	if uploadDate != "" {
		uploadDateFormatted := NewDate(uploadDate)
		upload.UploadDate = uploadDateFormatted
	}

	return upload, nil
}
