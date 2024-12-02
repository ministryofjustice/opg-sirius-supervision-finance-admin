package reports

import (
	"context"
	"encoding/csv"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin-api/db"
	"os"
)

type DbClient interface {
	Run(ctx context.Context, query db.ReportQuery) ([][]string, error)
}

type Client struct {
	db DbClient
}

func NewClient(ctx context.Context) (*Client, error) {
	dbClient, err := db.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer dbClient.Close(ctx)

	return &Client{db: dbClient}, nil
}

func (c *Client) Generate(ctx context.Context, filename string, query db.ReportQuery) (*os.File, error) {
	rows, err := c.db.Run(ctx, query)
	if err != nil {
		return nil, err
	}

	return createCsv(filename, rows)
}

func createCsv(filename string, items [][]string) (*os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	writer := csv.NewWriter(file)

	for _, item := range items {
		err = writer.Write(item)
		if err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if writer.Error() != nil {
		return nil, writer.Error()
	}

	return os.Open(filename)
}
