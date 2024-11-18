package api

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/smithy-go"
	"github.com/opg-sirius-finance-admin/apierror"
	"github.com/opg-sirius-finance-admin/shared"
	"io"
	"net/http"
	"os"
)

func (s *Server) download(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	var download shared.Download
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&download); err != nil {
		return err
	}

	c, err := os.ReadFile("/app/db/queries/aged_debt.sql")
	if err != nil {
		return err
	}
	sql := string(c)

	rows, err := s.conn.Query(ctx, sql)

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

	f, err := os.Create("test.csv")
	if err != nil {
		return err
	}
	writer := csv.NewWriter(f)

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

	c, err = os.ReadFile("test.csv")
	if err != nil {
		return err
	}

	err = s.filestorage.PutFile(
		ctx,
		os.Getenv("ASYNC_S3_BUCKET"),
		fmt.Sprintf("%s/%s", s3Directory, "test.csv"),
		bytes.NewReader(c))

	if err != nil {
		return err
	}

	payload := NotifyPayload{
		EmailAddress: "test@email.com",
		TemplateId:   "8c85cf6c-695f-493a-a25f-77b4fb5f6a8e",
		Personalisation: struct {
			Filename string `json:"upload_type"`
		}{
			"test",
		},
	}

	err = s.SendEmailToNotify(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}
