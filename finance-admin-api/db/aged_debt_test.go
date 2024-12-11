package db

import (
	fh "github.com/ministryofjustice/opg-sirius-supervision-finance-hub/shared"
	"github.com/stretchr/testify/assert"
)

func (suite *IntegrationSuite) Test_aged_debt() {
	ctx := suite.ctx
	oneHundred := "100.00"
	date := "2024-01-01"
	client1ID := suite.seeder.CreateClient(ctx, "Ian", "Test", "12345678")
	suite.seeder.CreateDeputy(ctx, client1ID, "Suzie", "Deputy", "LAY")
	suite.seeder.CreateInvoice(ctx, client1ID, fh.InvoiceTypeAD, &oneHundred, &date, nil, nil, nil)
	paidInvoiceID := suite.seeder.CreateInvoice(ctx, client1ID, fh.InvoiceTypeAD, &oneHundred, &date, nil, nil, nil)
	writeOffID := suite.seeder.CreateAdjustment(ctx, client1ID, paidInvoiceID, fh.AdjustmentTypeWriteOff, 100.00, "Written off")
	suite.seeder.ApproveAdjustment(ctx, client1ID, writeOffID)

	c := Client{suite.seeder.Conn}

	rows, err := c.Run(ctx, &AgedDebt{})
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, len(rows))

	results := mapByHeader(rows)
	assert.NotEmpty(suite.T(), results)

	assert.Equal(suite.T(), "Ian Test", results[0]["Customer Name"], "Customer Name")
	assert.Equal(suite.T(), "12345678", results[0]["Customer number"], "Customer number")
	assert.Equal(suite.T(), "100.00", results[0]["Original amount"], "Original amount")
	assert.Equal(suite.T(), "100.00", results[0]["Outstanding amount"], "Outstanding amount")
	//assert.Equal(suite.T(), "SOP number", results[0]["SOP number"], "SOP number")
	assert.Equal(suite.T(), "LAY", results[0]["Deputy type"], "Deputy type")
	//assert.Equal(suite.T(), "Active case?", results[0]["Active case?"], "Active case?")
	assert.Equal(suite.T(), "=\"0470\"", results[0]["Entity"], "Entity")
	assert.Equal(suite.T(), "99999999", results[0]["Receivable cost centre"], "Receivable cost centre")
	assert.Equal(suite.T(), "BALANCE SHEET", results[0]["Receivable cost centre description"], "Receivable cost centre description")
	assert.Equal(suite.T(), "1816100000", results[0]["Receivable account code"], "Receivable account code")
	assert.Equal(suite.T(), "10482009", results[0]["Revenue cost centre"], "Revenue cost centre")
	assert.Equal(suite.T(), "Supervision Investigations", results[0]["Revenue cost centre description"], "Revenue cost centre description")
	assert.Equal(suite.T(), "4481102093", results[0]["Revenue account code"], "Revenue account code")
	assert.Equal(suite.T(), "INC - RECEIPT OF FEES AND CHARGES - Appoint Deputy", results[0]["Revenue account code description"], "Revenue account code description")
	assert.Equal(suite.T(), "AD", results[0]["Invoice type"], "Invoice type")
	//assert.Equal(suite.T(), "Trx number", results[0]["Trx number"], "Trx number")
	assert.Equal(suite.T(), "AD - Assessment deputy invoice", results[0]["Transaction Description"], "Transaction Description")
	assert.Equal(suite.T(), "2024-01-01", results[0]["Invoice date"], "Invoice date")
	assert.Equal(suite.T(), "2024-01-31", results[0]["Due date"], "Due date")
	assert.Equal(suite.T(), "2023/24", results[0]["Financial year"], "Financial year")
	assert.Equal(suite.T(), "30 NET", results[0]["Payment terms"], "Payment terms")
	assert.Equal(suite.T(), "100.00", results[0]["Original amount"], "Original amount")
	assert.Equal(suite.T(), "100.00", results[0]["Outstanding amount"], "Outstanding amount")
	assert.Equal(suite.T(), "0", results[0]["Current"], "Current")
	assert.Equal(suite.T(), "100.00", results[0]["0-1 years"], "0-1 years")
	assert.Equal(suite.T(), "0", results[0]["1-2 years"], "1-2 years")
	assert.Equal(suite.T(), "0", results[0]["2-3 years"], "2-3 years")
	assert.Equal(suite.T(), "0", results[0]["3-5 years"], "3-5 years")
	assert.Equal(suite.T(), "0", results[0]["5+ years"], "5+ years")
	assert.Equal(suite.T(), "=\"0-1\"", results[0]["Debt impairment years"], "Debt impairment years")
}
