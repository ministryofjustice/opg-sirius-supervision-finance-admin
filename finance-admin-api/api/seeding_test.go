package api

import "github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin-api/testhelpers"

func (suite *IntegrationSuite) Test_seeding() {
	ctx := suite.ctx
	client := suite.seeder.CreateClient(ctx, &testhelpers.FinanceClient{
		Person: testhelpers.Person{
			FirstName: "Ian",
			Surname:   "Test",
		},
	})
	_, err := suite.seeder.CreateInvoice(ctx, client.Id, map[string]interface{}{})

	if err == nil {
		suite.T().Error("Seeding failed")
	}
	suite.T().Log("Seeding successful")
}
