package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/ministryofjustice/opg-go-common/env"
	"net/url"
)

type DBClient interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Close(ctx context.Context) error
}

type Client struct {
	db DBClient
}

func NewClient(ctx context.Context) (*Client, error) {
	dbConn := env.Get("POSTGRES_CONN", "")
	dbUser := env.Get("POSTGRES_USER", "")
	dbPassword := env.Get("POSTGRES_PASSWORD", "")
	pgDb := env.Get("POSTGRES_DB", "")

	conn, err := pgx.Connect(ctx, fmt.Sprintf("postgresql://%s:%s@%s/%s", dbUser, url.QueryEscape(dbPassword), dbConn, pgDb))

	if err != nil {
		return nil, err
	}

	return &Client{conn}, nil
}

func (c *Client) Close(ctx context.Context) error {
	return c.db.Close(ctx)
}

type ReportQuery interface {
	GetHeaders() []string
	GetQuery() string
	GetParams() []string
}

func (c *Client) Run(ctx context.Context, query ReportQuery) ([][]string, error) {
	items := [][]string{query.GetHeaders()}

	rows, err := c.db.Query(ctx, AgedDebtQuery, query.GetParams())

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
