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
			dateString:   "2020-01-02",
			wantErr:      false,
			wantFilename: "",
		},
		{
			name:         "Moto card payments report type",
			uploadType:   ReportTypeUploadPaymentsMOTOCard,
			dateString:   "2020-01-02",
			wantErr:      false,
			wantFilename: "feemoto_02:01:2020normal.csv",
		},
		{
			name:         "Online card payments report type",
			uploadType:   ReportTypeUploadPaymentsOnlineCard,
			dateString:   "2024-12-03",
			wantErr:      false,
			wantFilename: "feemoto_03:12:2024mlpayments.csv",
		},
		{
			name:         "Invalid date",
			uploadType:   ReportTypeUploadPaymentsMOTOCard,
			dateString:   "02/01/2020",
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
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
