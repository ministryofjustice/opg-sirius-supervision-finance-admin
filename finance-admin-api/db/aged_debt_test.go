package db

import (
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin-api/db/testhelpers"
	fh "github.com/ministryofjustice/opg-sirius-supervision-finance-hub/shared"
	"github.com/stretchr/testify/assert"
)

func (suite *IntegrationSuite) Test_aged_debt() {
	ctx := suite.ctx
	fc := suite.seeder.CreateClient(ctx, &testhelpers.FinanceClient{
		CourtRef: "12345678",
		Person: testhelpers.Person{
			FirstName: "Ian",
			Surname:   "Test",
		},
	})
	inv := fh.AddManualInvoice{
		InvoiceType: fh.InvoiceTypeAD,
		Amount: fh.Nillable[int]{
			Value: 10000,
			Valid: true,
		},
		RaisedDate: fh.Nillable[fh.Date]{
			Value: fh.NewDate("2024-01-01"),
			Valid: true,
		},
		StartDate: fh.Nillable[fh.Date]{
			Value: fh.NewDate("2024-01-01"),
			Valid: true,
		},
		EndDate: fh.Nillable[fh.Date]{
			Value: fh.NewDate("2024-01-01"),
			Valid: true,
		},
	}
	_ = suite.seeder.CreateInvoice(ctx, fc.Id, inv)

	c := Client{suite.seeder.Conn}

	rows, err := c.Run(ctx, &AgedDebt{})
	assert.NoError(suite.T(), err)

	results := mapByHeader(rows)
	assert.NotEmpty(suite.T(), results)

	assert.Equal(suite.T(), "Ian Test", results[0]["Customer Name"])
	assert.Equal(suite.T(), "12345678", results[0]["Customer number"])
	assert.Equal(suite.T(), "100.00", results[0]["Original amount"])
	assert.Equal(suite.T(), "100.00", results[0]["Outstanding amount"])
}
