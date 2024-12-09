package testhelpers

import (
	"context"
	fh "github.com/ministryofjustice/opg-sirius-supervision-finance-hub/shared"
	"log"
	"strconv"
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
	err := s.Conn.QueryRow(ctx, "INSERT INTO public.persons VALUES (NEXTVAL('public.persons_id_seq'), $1, $2, $3) RETURNING id", data.FirstName, data.Surname, data.CourtRef).Scan(&data.Id)
	if err != nil {
		log.Fatalf("failed to add FinanceClient: %v", err)
	}
	err = s.Conn.QueryRow(ctx, "INSERT INTO supervision_finance.finance_client VALUES (NEXTVAL('supervision_finance.finance_client_id_seq'), $1, '', 'DEMANDED') RETURNING id", data.Id).Scan(&data.FinanceClientId)
	if err != nil {
		log.Fatalf("failed to add finance_client: %v", err)
	}
	return data
}

func (s *Seeder) CreateDeputy(ctx context.Context, data *Deputy) int {
	var deputyId int
	err := s.Conn.QueryRow(ctx, "INSERT INTO public.persons VALUES (NEXTVAL('public.persons_id_seq'), $1, $2, NULL, $3, $4) RETURNING id", data.FirstName, data.Surname, data.client.Id, data.deputyType).Scan(&deputyId)
	if err != nil {
		log.Fatalf("failed to add Deputy: %v", err)
	}
	_, err = s.Conn.Exec(ctx, "UPDATE public.persons SET feepayer_id = $1 WHERE id = $2", deputyId, data.client.Id)
	if err != nil {
		log.Fatalf("failed to add Deputy to FinanceClient: %v", err)
	}
	return deputyId
}

func (s *Seeder) CreateInvoice(ctx context.Context, clientID int, data fh.AddManualInvoice) int {
	res, err := s.SendDataToAPI(ctx, "clients/"+strconv.Itoa(clientID)+"/invoices", data)
	if err != nil {
		log.Fatalf("failed to add invoice: %v", err)
	}
	var id int
	if res.StatusCode != 201 {
		log.Fatalf("failed to add invoice: status %v", res.Status)
	}
	err = s.Conn.QueryRow(ctx, "SELECT id FROM supervision_finance.invoice ORDER BY id DESC LIMIT 1").Scan(&id)
	if err != nil {
		log.Fatalf("failed find created invoice: %v", err)
	}
	return id
}
