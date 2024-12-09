package api

import (
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin-api/testhelpers"
	fh "github.com/ministryofjustice/opg-sirius-supervision-finance-hub/shared"
	"github.com/stretchr/testify/assert"
)

func (suite *IntegrationSuite) Test_seeding() {
	ctx := suite.ctx
	client := suite.seeder.CreateClient(ctx, &testhelpers.FinanceClient{
		CourtRef: "12345678",
		Person: testhelpers.Person{
			FirstName: "Ian",
			Surname:   "Test",
		},
	})
	invoice := fh.AddManualInvoice{
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
	invoiceId := suite.seeder.CreateInvoice(ctx, client.Id, invoice)
	assert.Equal(suite.T(), 1, invoiceId)
	suite.T().Log("Seeding successful")
}
