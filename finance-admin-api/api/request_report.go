package api

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/opg-sirius-finance-admin/db"
	"github.com/opg-sirius-finance-admin/shared"
	"net/http"
	"os"
	"time"
)

func (s *Server) requestReport(w http.ResponseWriter, r *http.Request) error {
	requestedDate := time.Now()
	ctx := context.Background()

	var download shared.Download
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&download); err != nil {
		return err
	}

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

	ef, err := os.Create("test.csv")
	if err != nil {
		return err
	}

	defer ef.Close()

	query := db.AgedDebtQuery

	rows, err := s.conn.Query(ctx, query)
	if err != nil {
		return err
	}

	defer rows.Close()
	var items [][]string
	for rows.Next() {
		var i []string
		var stringValue string

		values, err := rows.Values()
		if err != nil {
			return err
		}
		for _, value := range values {
			stringValue, _ = value.(string)
			i = append(i, stringValue)
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	writer := csv.NewWriter(ef)

	err = writer.Write(agedDebtHeaders)
	if err != nil {
		return err
	}

	for _, item := range items {
		err = writer.Write(item)
		if err != nil {
			return err
		}
	}

	writer.Flush()
	if writer.Error() != nil {
		return writer.Error()
	}

	ef.Close()

	rf, err := os.Open("test.csv")
	if err != nil {
		return err
	}

	versionId, err := s.filestorage.PutFile(
		ctx,
		os.Getenv("REPORTS_S3_BUCKET"),
		"test.csv",
		rf,
	)

	rf.Close()

	if err != nil {
		return err
	}

	if versionId == nil {
		return fmt.Errorf("S3 version ID not found")
	}

	downloadRequest := shared.DownloadRequest{
		Key:       "test.csv",
		VersionId: *versionId,
	}

	uid, err := downloadRequest.Encode()
	if err != nil {
		return err
	}

	siriusUrl := os.Getenv("SIRIUS_PUBLIC_URL")
	prefix := os.Getenv("PREFIX")
	downloadLink := siriusUrl + prefix + "/download?uid=" + uid

	// JSON.stringify({ Key: key, VersionId: versionId })

	payload := NotifyPayload{
		EmailAddress: "test@email.com",
		TemplateId:   "bade69e4-0eb1-4896-a709-bd8f8371a629",
		Personalisation: struct {
			FileLink          string `json:"file_link"`
			ReportName        string `json:"report_name"`
			RequestedDate     string `json:"requested_date"`
			RequestedDateTime string `json:"requested_date_time"`
		}{
			downloadLink,
			"Aged Debt",
			requestedDate.Format("2006-01-02"),
			requestedDate.Format("2006-01-02 15:04:05"),
		},
	}

	err = s.SendEmailToNotify(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}
