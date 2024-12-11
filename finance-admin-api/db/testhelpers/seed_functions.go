package testhelpers

import (
	"context"
	"fmt"
	fh "github.com/ministryofjustice/opg-sirius-supervision-finance-hub/shared"
	"io"
	"log"
	"net/http"
	"strconv"
)

func (s *Seeder) CreateClient(ctx context.Context, firstName string, surname string, courtRef string, sopNumber string) int {
	var (
		clientId        int
		financeClientId int
	)
	err := s.Conn.QueryRow(ctx, "INSERT INTO public.persons VALUES (NEXTVAL('public.persons_id_seq'), $1, $2, $3) RETURNING id", firstName, surname, courtRef).Scan(&clientId)
	if err != nil {
		log.Fatalf("failed to add FinanceClient: %v", err)
	}
	err = s.Conn.QueryRow(ctx, "INSERT INTO supervision_finance.finance_client VALUES (NEXTVAL('supervision_finance.finance_client_id_seq'), $1, $2, 'DEMANDED') RETURNING id", clientId, sopNumber).Scan(&financeClientId)
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

func (s *Seeder) CreateOrder(ctx context.Context, clientId int, status string) {
	var (
		personId int
	)
	err := s.Conn.QueryRow(ctx, "SELECT client_id FROM supervision_finance.finance_client WHERE id = $1", clientId).Scan(&personId)
	if err != nil {
		log.Fatalf("failed find client record: %v", err)
	}
	_, err = s.Conn.Exec(ctx, "INSERT INTO public.cases VALUES (NEXTVAL('public.cases_id_seq'), $1, $2)", personId, status)
	if err != nil {
		log.Fatalf("failed to add order: %v", err)
	}
}

func (s *Seeder) CreateInvoice(ctx context.Context, clientID int, invoiceType fh.InvoiceType, amount *string, raisedDate *string, startDate *string, endDate *string, supervisionLevel *string) (int, string) {
	invoice := fh.AddManualInvoice{
		InvoiceType:      invoiceType,
		Amount:           fh.TransformNillableInt(amount),
		RaisedDate:       fh.TransformNillableDate(raisedDate),
		StartDate:        fh.TransformNillableDate(startDate),
		EndDate:          fh.TransformNillableDate(endDate),
		SupervisionLevel: fh.TransformNillableString(supervisionLevel),
	}

	res, _ := s.SendDataToAPI(ctx, http.MethodPost, "clients/"+strconv.Itoa(clientID)+"/invoices", invoice)
	var (
		id        int
		reference string
	)
	if res.StatusCode != 201 {
		body, _ := io.ReadAll(res.Body)
		log.Fatalf("failed to add invoice: status %v, body: %v", res.Status, string(body))

	}
	err := s.Conn.QueryRow(ctx, "SELECT id, reference FROM supervision_finance.invoice ORDER BY id DESC LIMIT 1").Scan(&id, &reference)
	if err != nil {
		log.Fatalf("failed find created invoice: %v", err)
	}
	return id, reference
}

func (s *Seeder) CreateAdjustment(ctx context.Context, clientID int, invoiceId int, adjustmentType fh.AdjustmentType, amount int, notes string) int {
	adjustment := fh.AddInvoiceAdjustmentRequest{
		AdjustmentType:  adjustmentType,
		AdjustmentNotes: notes,
		Amount:          amount,
	}
	res, _ := s.SendDataToAPI(ctx, http.MethodPost, fmt.Sprintf("clients/%d/invoices/%d/invoice-adjustments", clientID, invoiceId), adjustment)
	var id int
	if res.StatusCode != 201 {
		log.Fatalf("failed to add adjustment: status %v", res.Status)
	}
	err := s.Conn.QueryRow(ctx, "SELECT id FROM supervision_finance.invoice_adjustment ORDER BY id DESC LIMIT 1").Scan(&id)
	if err != nil {
		log.Fatalf("failed find created adjustment: %v", err)
	}
	return id
}

func (s *Seeder) ApproveAdjustment(ctx context.Context, clientID int, adjustmentId int) {
	decision := fh.UpdateInvoiceAdjustment{
		Status: fh.AdjustmentStatusApproved,
	}
	res, _ := s.SendDataToAPI(ctx, http.MethodPut, fmt.Sprintf("clients/%d/invoice-adjustments/%d", clientID, adjustmentId), decision)
	if res.StatusCode != 204 {
		log.Fatalf("failed to approve adjustment: status %v", res.Status)
	}
}
