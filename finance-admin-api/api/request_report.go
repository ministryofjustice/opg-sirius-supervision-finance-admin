package api

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/opg-sirius-finance-admin/shared"
	"net/http"
	"os"
)

func (s *Server) requestReport(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	var download shared.Download
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&download); err != nil {
		return err
	}

	c, err := os.ReadFile("../db/queries/aged_debt.sql")
	if err != nil {
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

	query := string(c)

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

	versionId, err := s.filestorage.PutFile(
		ctx,
		os.Getenv("ASYNC_S3_BUCKET"),
		fmt.Sprintf("%s/%s", s3Directory, "test.csv"),
		ef,
	)

	if err != nil {
		return err
	}

	if versionId == nil {
		return fmt.Errorf("S3 version ID not found")
	}

	downloadRequest := shared.DownloadRequest{
		Key:       fmt.Sprintf("%s/%s", s3Directory, "test.csv"),
		VersionId: *versionId,
	}

	uid, err := downloadRequest.Encode()
	if err != nil {
		return err
	}

	payload := NotifyPayload{
		EmailAddress: "test@email.com",
		TemplateId:   "8c85cf6c-695f-493a-a25f-77b4fb5f6a8e",
		Personalisation: struct {
			Uid string `json:"upload_type"`
		}{
			uid,
		},
	}

	err = s.SendEmailToNotify(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}
