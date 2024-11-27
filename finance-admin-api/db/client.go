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
	GetParams() []any
}

func (c *Client) Run(ctx context.Context, query ReportQuery) ([][]string, error) {
	headers := [][]string{query.GetHeaders()}

	rows, err := c.db.Query(ctx, query.GetQuery(), query.GetParams()...)
	if err != nil {
		return nil, err
	}

	stringRows, err := pgx.CollectRows[[]string](rows, rowToStringMap)

	if err != nil {
		return nil, err
	}

	return append(headers, stringRows...), nil
}

func rowToStringMap(row pgx.CollectableRow) ([]string, error) {
	var stringRow []string
	values, err := row.Values()
	if err != nil {
		return nil, err
	}

	for _, value := range values {
		stringRow = append(stringRow, fmt.Sprintf("%v", value))
	}
	return stringRow, nil
}
