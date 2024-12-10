package db

import (
	fh "github.com/ministryofjustice/opg-sirius-supervision-finance-hub/shared"
	"github.com/stretchr/testify/assert"
)

func (suite *IntegrationSuite) Test_aged_debt() {
	ctx := suite.ctx
	client1ID := suite.seeder.CreateClient(ctx, "Ian", "Test", "12345678")
	suite.seeder.CreateDeputy(ctx, client1ID, "Suzie", "Deputy", "LAY")
	suite.seeder.CreateInvoice(ctx, client1ID, fh.InvoiceTypeAD, "100.00", "2024-01-01", "2024-01-01", "2024-01-01", "")
	paidInvoiceID := suite.seeder.CreateInvoice(ctx, client1ID, fh.InvoiceTypeAD, "100.00", "2024-01-01", "", "", "")
	writeOffID := suite.seeder.CreateAdjustment(ctx, client1ID, paidInvoiceID, fh.AdjustmentTypeWriteOff, 100.00, "Written off")
	suite.seeder.ApproveAdjustment(ctx, client1ID, paidInvoiceID, writeOffID)

	c := Client{suite.seeder.Conn}

	rows, err := c.Run(ctx, &AgedDebt{})
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, len(rows))

	results := mapByHeader(rows)
	assert.NotEmpty(suite.T(), results)

	assert.Equal(suite.T(), "Ian Test", results[0]["Customer Name"])
	assert.Equal(suite.T(), "12345678", results[0]["Customer number"])
	assert.Equal(suite.T(), "100.00", results[0]["Original amount"])
	assert.Equal(suite.T(), "100.00", results[0]["Outstanding amount"])
}
