package api

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/opg-sirius-finance-admin/apierror"
	"github.com/opg-sirius-finance-admin/db"
	"github.com/opg-sirius-finance-admin/shared"
	"net/http"
	"os"
	"time"
)

const reportRequestedTemplateId = "872d88b3-076e-495c-bf81-a2be2d3d234c"

func (s *Server) requestReport(w http.ResponseWriter, r *http.Request) error {
	var download shared.Download
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&download); err != nil {
		return err
	}

	if download.Email == "" {
		return apierror.ValidationError{Errors: apierror.ValidationErrors{
			"Email": {
				"required": "This field Email needs to be looked at required",
			},
		},
		}
	}

	go func() {
		err := s.generateAndUploadReport(context.Background(), download, time.Now())
		if err != nil {
			fmt.Println(err)
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	return nil
}

func (s *Server) generateAndUploadReport(ctx context.Context, download shared.Download, requestedDate time.Time) error {
	var items [][]string
	var filename string
	var reportName string
	var err error

	switch download.ReportAccountType {
	case "AgedDebt":
		//parsedDate, err := time.Parse("02/01/2006", requestedDate)
		//if err != nil {
		//	return err
		//}
		filename = fmt.Sprintf("ageddebt_%s.csv", requestedDate.Format("02:01:2006"))
		reportName = "Aged Debt"
		items, err = s.requestAgedDebtReport(ctx)
		if err != nil {
			return err
		}
	}

	file, err := createCsv(filename, items)
	if err != nil {
		return err
	}

	versionId, err := s.filestorage.PutFile(
		ctx,
		os.Getenv("REPORTS_S3_BUCKET"),
		filename,
		file,
	)

	file.Close()

	if err != nil {
		return err
	}

	payload, err := createDownloadNotifyPayload(download.Email, filename, versionId, requestedDate, reportName)
	if err != nil {
		return err
	}

	err = s.SendEmailToNotify(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) requestAgedDebtReport(ctx context.Context) ([][]string, error) {
	agedDebtHeaders := []string{
		"Customer Name",
		"Customer number",
		"SOP number",
		"Deputy type",
		"Active case?",
		"Entity",
		"Receivable cost centre",
		"Receivable cost centre description",
		"Receivable account code",
		"Revenue cost centre",
		"Revenue cost centre description",
		"Revenue account code",
		"Revenue account code description",
		"Invoice type",
		"Trx number",
		"Transaction Description",
		"Invoice date",
		"Due date",
		"Financial year",
		"Payment terms",
		"Original amount",
		"Outstanding amount",
		"Current",
		"0-1 years",
		"1-2 years",
		"2-3 years",
		"3-5 years",
		"5+ years",
		"Debt impairment years",
	}

	items := [][]string{agedDebtHeaders}

	rows, err := s.conn.Query(ctx, db.AgedDebtQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var i []string
		var stringValue string

		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		for _, value := range values {
			stringValue, _ = value.(string)
			i = append(i, stringValue)
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
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

	file.Close()

	rf, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return rf, nil
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

	siriusUrl := os.Getenv("SIRIUS_PUBLIC_URL")
	prefix := os.Getenv("PREFIX")
	downloadLink := siriusUrl + prefix + "/download?uid=" + uid

	payload := NotifyPayload{
		EmailAddress: emailAddress,
		TemplateId:   reportRequestedTemplateId,
		Personalisation: struct {
			FileLink          string `json:"file_link"`
			ReportName        string `json:"report_name"`
			RequestedDate     string `json:"requested_date"`
			RequestedDateTime string `json:"requested_date_time"`
		}{
			downloadLink,
			reportName,
			requestedDate.Format("2006-01-02"),
			requestedDate.Format("2006-01-02 15:04:05"),
		},
	}

	return payload, nil
}
