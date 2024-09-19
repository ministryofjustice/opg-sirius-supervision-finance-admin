package model

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

type errorReader struct{}

func (e *errorReader) Read(_ []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func TestNewUploadReturnsCorrectly(t *testing.T) {
	reportUploadType := "SomeReportType"
	email := "test@example.com"
	uploadDate := "11/05/2024"
	fileContent := []byte("file content")
	fileReader := bytes.NewReader(fileContent)

	// Expected result
	expectedDate := NewDate(uploadDate)
	expectedUpload := Upload{
		ReportUploadType: reportUploadType,
		Email:            email,
		File:             fileContent,
		UploadDate:       &expectedDate,
	}

	upload, err := NewUpload(reportUploadType, uploadDate, email, fileReader)

	assert.NoError(t, err)
	assert.Equal(t, expectedUpload, upload)
}

func TestNewUploadWithNoFileReturnsCorrectly(t *testing.T) {
	reportUploadType := "SomeReportType"
	email := "test@example.com"
	uploadDate := ""
	fileContent := []byte("file content")
	fileReader := bytes.NewReader(fileContent)

	expectedUpload := Upload{
		ReportUploadType: reportUploadType,
		Email:            email,
		File:             fileContent,
		UploadDate:       nil,
	}

	upload, err := NewUpload(reportUploadType, uploadDate, email, fileReader)

	assert.NoError(t, err)
	assert.Equal(t, expectedUpload, upload)
}

func TestNewUpload_ReturnsError(t *testing.T) {
	// Prepare input with an errorReader that always fails
	reportUploadType := "SomeReportType"
	email := "test@example.com"
	uploadDate := "11/05/2024"
	reader := &errorReader{}

	upload, err := NewUpload(reportUploadType, uploadDate, email, reader)

	assert.ErrorIs(t, err, io.ErrUnexpectedEOF)
	assert.Empty(t, upload.File)
}
