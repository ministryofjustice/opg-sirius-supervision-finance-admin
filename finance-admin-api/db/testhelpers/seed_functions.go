package testhelpers

import (
	"context"
	"fmt"
	fh "github.com/ministryofjustice/opg-sirius-supervision-finance-hub/shared"
	"log"
	"net/http"
	"strconv"
)

func (s *Seeder) CreateClient(ctx context.Context, firstName string, surname string, courtRef string) int {
	var (
		clientId        int
		financeClientId int
	)
	err := s.Conn.QueryRow(ctx, "INSERT INTO public.persons VALUES (NEXTVAL('public.persons_id_seq'), $1, $2, $3) RETURNING id", firstName, surname, courtRef).Scan(&clientId)
	if err != nil {
		log.Fatalf("failed to add FinanceClient: %v", err)
	}
	err = s.Conn.QueryRow(ctx, "INSERT INTO supervision_finance.finance_client VALUES (NEXTVAL('supervision_finance.finance_client_id_seq'), $1, '', 'DEMANDED') RETURNING id", clientId).Scan(&financeClientId)
	if err != nil {
		log.Fatalf("failed to add finance_client: %v", err)
	}
	return financeClientId
}

func (s *Seeder) CreateDeputy(ctx context.Context, clientId int, firstName string, surname string, deputyType string) int {
	var (
		personId int
		deputyId int
	)
	err := s.Conn.QueryRow(ctx, "SELECT client_id FROM supervision_finance.finance_client WHERE id = $1", clientId).Scan(&personId)
	if err != nil {
		log.Fatalf("failed find client record: %v", err)
	}
	err = s.Conn.QueryRow(ctx, "INSERT INTO public.persons VALUES (NEXTVAL('public.persons_id_seq'), $1, $2, NULL, $3, $4) RETURNING id", firstName, surname, personId, deputyType).Scan(&deputyId)
	if err != nil {
		log.Fatalf("failed to add Deputy: %v", err)
	}
	_, err = s.Conn.Exec(ctx, "UPDATE public.persons SET feepayer_id = $1 WHERE id = $2", deputyId, personId)
	if err != nil {
		log.Fatalf("failed to add Deputy to FinanceClient: %v", err)
	}
	return deputyId
}

func (s *Seeder) CreateInvoice(ctx context.Context, clientID int, invoiceType fh.InvoiceType, amount string, raisedDate string, startDate string, endDate string, supervisionLevel string) int {
	invoice := fh.AddManualInvoice{
		InvoiceType:      invoiceType,
		Amount:           fh.TransformNillableInt(&amount),
		RaisedDate:       fh.TransformNillableDate(&raisedDate),
		StartDate:        fh.TransformNillableDate(&startDate),
		EndDate:          fh.TransformNillableDate(&endDate),
		SupervisionLevel: fh.TransformNillableString(&supervisionLevel),
	}

	res, err := s.SendDataToAPI(ctx, http.MethodPost, "clients/"+strconv.Itoa(clientID)+"/invoices", invoice)
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

func (s *Seeder) CreateAdjustment(ctx context.Context, clientID int, invoiceId int, adjustmentType fh.AdjustmentType, amount int, notes string) int {
	adjustment := fh.AddInvoiceAdjustmentRequest{
		AdjustmentType:  adjustmentType,
		AdjustmentNotes: notes,
		Amount:          amount,
	}
	res, err := s.SendDataToAPI(ctx, http.MethodPost, fmt.Sprintf("/clients/%d/invoices/%d/invoice-adjustments", clientID, invoiceId), adjustment)
	if err != nil {
		log.Fatalf("failed to add adjustment: %v", err)
	}
	var id int
	if res.StatusCode != 201 {
		log.Fatalf("failed to add adjustment: status %v", res.Status)
	}
	err = s.Conn.QueryRow(ctx, "SELECT id FROM supervision_finance.invoice_adjustment ORDER BY id DESC LIMIT 1").Scan(&id)
	if err != nil {
		log.Fatalf("failed find created adjustment: %v", err)
	}
	return id
}

func (s *Seeder) ApproveAdjustment(ctx context.Context, clientID int, invoiceId int, adjustmentId int) {
	decision := fh.UpdateInvoiceAdjustment{
		Status: fh.AdjustmentStatusApproved,
	}
	res, err := s.SendDataToAPI(ctx, http.MethodPut, fmt.Sprintf("/clients/%d/invoices/%d/invoice-adjustments/%d", clientID, invoiceId, adjustmentId), decision)
	if err != nil {
		log.Fatalf("failed to approve adjustment: %v", err)
	}
	if res.StatusCode != 200 {
		log.Fatalf("failed to approve adjustment: status %v", res.Status)
	}
}
