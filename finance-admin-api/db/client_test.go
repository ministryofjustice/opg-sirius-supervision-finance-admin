package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockRow struct {
	values       [][]any
	err          error
	recordNumber int
}

func (m *mockRow) Close()                                       {}
func (m *mockRow) Err() error                                   { return nil }
func (m *mockRow) CommandTag() pgconn.CommandTag                { return pgconn.CommandTag{} }
func (m *mockRow) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (m *mockRow) Scan(dest ...any) error                       { return nil }
func (m *mockRow) RawValues() [][]byte                          { return nil }
func (m *mockRow) Conn() *pgx.Conn                              { return nil }

func (m *mockRow) Next() bool {
	m.recordNumber++
	return m.recordNumber <= len(m.values)
}

func (m *mockRow) Values() ([]any, error) {
	recordNumber := m.recordNumber
	if recordNumber != 0 {
		recordNumber--
	}
	return m.values[recordNumber], m.err
}

type mockDbClient struct {
	values [][]any
	err    error
}

func (m mockDbClient) Close(ctx context.Context) error { return nil }

func (m mockDbClient) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	rows := mockRow{values: m.values}
	return &rows, m.err
}

func TestRowToStringMap(t *testing.T) {
	row := mockRow{
		values: [][]any{{"test", "values", 2, 3, 5.55}},
	}

	got, err := rowToStringMap(&row)

	want := []string{"test", "values", "2", "3", "5.55"}

	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

func TestRowToStringMapError(t *testing.T) {
	row := mockRow{
		values: [][]any{{"test", "values", 2, 3, 5.55}},
		err:    fmt.Errorf("Oh no!"),
	}

	got, err := rowToStringMap(&row)
	want := fmt.Errorf("Oh no!")

	assert.Nil(t, got)
	assert.Equal(t, want, err)
}

type mockQueryReport struct {
	headers []string
}

func (m mockQueryReport) GetQuery() string { return "" }
func (m mockQueryReport) GetParams() []any { return nil }

func (m mockQueryReport) GetHeaders() []string {
	return m.headers
}

func TestRun(t *testing.T) {
	values := [][]any{{"Joseph Smith", "123 Fake Street", 125}, {"Not Joseph Smith", "28 Real Avenue", 50000}}

	dbClient := mockDbClient{values: values}
	mockClient := Client{dbClient}
	ctx := context.Background()

	headers := []string{"Name", "Address", "Balance"}
	query := mockQueryReport{headers: headers}

	got, err := mockClient.Run(ctx, query)

	want := [][]string{
		{"Name", "Address", "Balance"},
		{"Joseph Smith", "123 Fake Street", "125"},
		{"Not Joseph Smith", "28 Real Avenue", "50000"},
	}

	assert.Equal(t, want, got)
	assert.Nil(t, err)
}

func TestRunError(t *testing.T) {
	dbClient := mockDbClient{err: fmt.Errorf("Oh dear!")}
	mockClient := Client{dbClient}
	ctx := context.Background()

	query := mockQueryReport{}

	got, err := mockClient.Run(ctx, query)

	want := fmt.Errorf("Oh dear!")

	assert.Equal(t, want, err)
	assert.Nil(t, got)
}
