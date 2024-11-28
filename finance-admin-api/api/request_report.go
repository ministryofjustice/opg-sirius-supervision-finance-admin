package api

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/apierror"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin-api/db"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/shared"
	"net/http"
	"os"
	"time"
)

const reportRequestedTemplateId = "bade69e4-0eb1-4896-a709-bd8f8371a629"

func (s *Server) requestReport(w http.ResponseWriter, r *http.Request) error {
	var reportRequest shared.ReportRequest
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&reportRequest); err != nil {
		return err
	}

	if reportRequest.Email == "" {
		return apierror.ValidationError{Errors: apierror.ValidationErrors{
			"Email": {
				"required": "This field Email needs to be looked at required",
			},
		},
		}
	}

	go func() {
		err := s.generateAndUploadReport(context.Background(), reportRequest, time.Now())
		if err != nil {
			logger := telemetry.LoggerFromContext(r.Context())
			logger.Error(err.Error())
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	return nil
}

func (s *Server) generateAndUploadReport(ctx context.Context, reportRequest shared.ReportRequest, requestedDate time.Time) error {
	var query db.ReportQuery
	var err error

	accountType := shared.ParseReportAccountType(reportRequest.ReportAccountType)
	filename := fmt.Sprintf("%s_%s.csv", accountType.Key(), requestedDate.Format("02:01:2006"))

	switch accountType {
	case shared.ReportAccountTypeAgedDebt:
		query = &db.AgedDebt{
			FromDate: reportRequest.FromDateField,
			ToDate:   reportRequest.ToDateField,
		}
	}

	rows, err := s.conn.Run(ctx, query)
	if err != nil {
		return err
	}

	file, err := createCsv(filename, rows)
	if err != nil {
		return err
	}

	defer file.Close()

	versionId, err := s.filestorage.PutFile(
		ctx,
		os.Getenv("REPORTS_S3_BUCKET"),
		filename,
		file,
	)

	if err != nil {
		return err
	}

	payload, err := createDownloadNotifyPayload(reportRequest.Email, filename, versionId, requestedDate, accountType.Translation())
	if err != nil {
		return err
	}

	err = s.SendEmailToNotify(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}

func createCsv(filename string, items [][]string) (*os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	writer := csv.NewWriter(file)

	for _, item := range items {
		err = writer.Write(item)
		if err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if writer.Error() != nil {
		return nil, writer.Error()
	}

	return os.Open(filename)
}

type reportRequestedNotifyPersonalisation struct {
	FileLink          string `json:"file_link"`
	ReportName        string `json:"report_name"`
	RequestedDate     string `json:"requested_date"`
	RequestedDateTime string `json:"requested_date_time"`
}

func createDownloadNotifyPayload(emailAddress string, filename string, versionId *string, requestedDate time.Time, reportName string) (NotifyPayload, error) {
	if versionId == nil {
		return NotifyPayload{}, fmt.Errorf("S3 version ID not found")
	}

	downloadRequest := shared.DownloadRequest{
		Key:       filename,
		VersionId: *versionId,
	}

	uid, err := downloadRequest.Encode()
	if err != nil {
		return NotifyPayload{}, err
	}

	downloadLink := fmt.Sprintf("%s%s/download?uid=%s", os.Getenv("SIRIUS_PUBLIC_URL"), os.Getenv("PREFIX"), uid)

	payload := NotifyPayload{
		EmailAddress: emailAddress,
		TemplateId:   reportRequestedTemplateId,
		Personalisation: reportRequestedNotifyPersonalisation{
			downloadLink,
			reportName,
			requestedDate.Format("2006-01-02"),
			requestedDate.Format("2006-01-02 15:04:05"),
		},
	}

	return payload, nil
}
