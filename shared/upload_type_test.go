package shared

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReportUploadType_Filename(t *testing.T) {
	tests := []struct {
		name         string
		uploadType   ReportUploadType
		dateString   string
		wantErr      bool
		wantFilename string
	}{
		{
			name:         "Non-moto card payments report type",
			uploadType:   ReportTypeUploadDeputySchedule,
			dateString:   "02/01/2020",
			wantErr:      false,
			wantFilename: "",
		},
		{
			name:         "Moto card payments report type",
			uploadType:   ReportTypeUploadPaymentsMOTOCard,
			dateString:   "02/01/2020",
			wantErr:      false,
			wantFilename: "",
		},
		{
			name:         "Invalid date",
			uploadType:   ReportTypeUploadPaymentsMOTOCard,
			dateString:   "hehe",
			wantErr:      true,
			wantFilename: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filename, err := test.uploadType.Filename(test.dateString)

			assert.Equal(t, test.wantFilename, filename)

			if test.wantErr {
				assert.Error(t, err)
			}
		})
	}
}
