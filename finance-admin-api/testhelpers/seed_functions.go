package testhelpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Person struct {
	FirstName string
	Surname   string
}

type FinanceClient struct {
	Id              int
	FinanceClientId int
	CourtRef        string
	Person
}

type Deputy struct {
	deputyType string
	client     *FinanceClient
	Person
}

func (s *Seeder) CreateClient(ctx context.Context, data *FinanceClient) *FinanceClient {
	err := s.GetConn().QueryRow(ctx, "INSERT INTO public.persons VALUES (NEXTVAL('public.persons_id_seq'), $1, $2, $3) RETURNING id", data.FirstName, data.Surname, data.CourtRef).Scan(data.Id)
	if err != nil {
		log.Fatalf("failed to add FinanceClient: %v", err)
	}
	err = s.GetConn().QueryRow(ctx, "INSERT INTO supervision_finance.finance_client VALUES (NEXTVAL('supervision_finance.finance_client_id_seq'), $1, '', 'DEMANDED') RETURNING id", data.Id).Scan(data.FinanceClientId)
	if err != nil {
		log.Fatalf("failed to add finance_client: %v", err)
	}
	return data
}

func (s *Seeder) CreateDeputy(ctx context.Context, data Deputy) int {
	var deputyId int
	err := s.GetConn().QueryRow(ctx, "INSERT INTO public.persons VALUES (NEXTVAL('public.persons_id_seq'), $1, $2, NULL, $3, $4) RETURNING id", data.FirstName, data.Surname, data.client.Id, data.deputyType).Scan(&deputyId)
	if err != nil {
		log.Fatalf("failed to add Deputy: %v", err)
	}
	_, err = s.GetConn().Exec(ctx, "UPDATE public.persons SET feepayer_id = $1 WHERE id = $2", deputyId, data.client.Id)
	if err != nil {
		log.Fatalf("failed to add Deputy to FinanceClient: %v", err)
	}
	return deputyId
}

func (s *Seeder) CreateInvoice(ctx context.Context, clientID int, invoiceData map[string]interface{}) (int, error) {
	url := fmt.Sprintf("%s/invoices", s.FinanceHub.BaseURL)

	data, err := json.Marshal(invoiceData)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.FinanceHub.HTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return 0, fmt.Errorf("failed to create invoice: %s", body)
	}

	var result struct {
		InvoiceID int `json:"invoice_id"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, err
	}

	return result.InvoiceID, nil
}
