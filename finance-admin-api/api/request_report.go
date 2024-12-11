package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/apierror"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin-api/db"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/shared"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const reportRequestedTemplateId = "bade69e4-0eb1-4896-a709-bd8f8371a629"

func validateReportRequest(reportRequest shared.ReportRequest) error {
	errors := apierror.ValidationErrors{}
	if reportRequest.Email == "" {
		errors["Email"] = map[string]string{"required": "This field Email needs to be looked at required"}
	}

	switch reportRequest.ReportAccountType {
	case shared.ReportAccountTypeBadDebtWriteOffReport, shared.ReportAccountTypePaidInvoiceReport:
		if reportRequest.FromDateField != nil {
			liveDate := shared.NewDate(os.Getenv("FINANCE_HUB_LIVE_DATE"))

			if reportRequest.FromDateField.Before(liveDate) {
				errors["FromDate"] = map[string]string{"date-after-live": "Date from cannot be before finance hub live date"}
			}
		}
	}

	if len(errors) > 0 {
		return apierror.ValidationError{Errors: errors}
	}

	return nil
}

func (s *Server) requestReport(w http.ResponseWriter, r *http.Request) error {
	var reportRequest shared.ReportRequest
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&reportRequest); err != nil {
		return err
	}

	if err := validateReportRequest(reportRequest); err != nil {
		return err
	}

	go func(logger *slog.Logger) {
		err := s.generateAndUploadReport(context.Background(), reportRequest, time.Now())
		if err != nil {
			logger.Error(err.Error())
		}
	}(telemetry.LoggerFromContext(r.Context()))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	return nil
}

func (s *Server) generateAndUploadReport(ctx context.Context, reportRequest shared.ReportRequest, requestedDate time.Time) error {
	var query db.ReportQuery
	var filename string
	var reportName string
	var err error

	switch reportRequest.ReportType {
	case shared.ReportsTypeAccountsReceivable:
		filename = fmt.Sprintf("%s_%s.csv", reportRequest.ReportAccountType.Key(), requestedDate.Format("02:01:2006"))
		reportName = reportRequest.ReportAccountType.Translation()

		switch reportRequest.ReportAccountType {
		case shared.ReportAccountTypeAgedDebt:
			query = &db.AgedDebt{
				FromDate: reportRequest.FromDateField,
				ToDate:   reportRequest.ToDateField,
			}
		case shared.ReportAccountTypeAgedDebtByCustomer:
			query = &db.AgedDebtByCustomer{}
		case shared.ReportAccountTypeBadDebtWriteOffReport:
			query = &db.BadDebtWriteOff{
				FromDate: reportRequest.FromDateField,
				ToDate:   reportRequest.ToDateField,
			}
		case shared.ReportAccountTypePaidInvoiceReport:
			query = &db.PaidInvoices{
				FromDate: reportRequest.FromDateField,
				ToDate:   reportRequest.ToDateField,
			}
		default:
			return fmt.Errorf("Unknown query")
		}
	}

	file, err := s.reports.Generate(ctx, filename, query)
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

	payload, err := createDownloadNotifyPayload(reportRequest.Email, filename, versionId, requestedDate, reportName)
	if err != nil {
		return err
	}

	err = s.SendEmailToNotify(ctx, payload)
	if err != nil {
		return err
	}

	return nil
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
