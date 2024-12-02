package testhelpers

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Client struct {
	db *pgxpool.Pool
}

type person struct {
	firstName string
	surname   string
}

type financeClient struct {
	id              int
	financeClientId int
	courtRef        string
	person
}

type deputy struct {
	deputyType string
	client     *financeClient
	person
}

// createClient creates a new financeClient in the public schema
func (p *Client) createClient(ctx context.Context, data *financeClient) *financeClient {
	err := p.db.QueryRow(ctx, "INSERT INTO public.persons VALUES (NEXTVAL('public.persons_id_seq'), $1, $2, $3) RETURNING id", data.firstName, data.surname, data.courtRef).Scan(data.id)
	if err != nil {
		log.Fatalf("failed to add financeClient: %v", err)
	}
	err = p.db.QueryRow(ctx, "INSERT INTO supervision_finance.finance_client VALUES (NEXTVAL('supervision_finance.finance_client_id_seq'), $1, '', 'DEMANDED') RETURNING id", data.id).Scan(data.financeClientId)
	if err != nil {
		log.Fatalf("failed to add finance_client: %v", err)
	}
	return data
}

// createDeputy creates a new deputy in the public schema
func (p *Client) createDeputy(ctx context.Context, data deputy) int {
	var deputyId int
	err := p.db.QueryRow(ctx, "INSERT INTO public.persons VALUES (NEXTVAL('public.persons_id_seq'), $1, $2, NULL, $3, $4) RETURNING id", data.firstName, data.surname, data.client.id, data.deputyType).Scan(&deputyId)
	if err != nil {
		log.Fatalf("failed to add deputy: %v", err)
	}
	_, err = p.db.Exec(ctx, "UPDATE public.persons SET feepayer_id = $1 WHERE id = $2", deputyId, data.client.id)
	if err != nil {
		log.Fatalf("failed to add deputy to financeClient: %v", err)
	}
	return deputyId
}

type order struct {
	clientId    int
	orderStatus string
}

// createOrder creates a new order in the public schema
func (p *Client) createOrder(ctx context.Context, data order) int {
	var orderId int
	err := p.db.QueryRow(ctx, "INSERT INTO public.cases VALUES (NEXTVAL('public.cases_id_seq'), $1, $2) RETURNING id", data.clientId, data.orderStatus).Scan(&orderId)
	if err != nil {
		log.Fatalf("failed to add order: %v", err)
	}
	return orderId
}
