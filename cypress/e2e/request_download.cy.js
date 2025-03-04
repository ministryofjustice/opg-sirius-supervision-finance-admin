describe("Downloads", () => {
    beforeEach(() => {
        cy.visit("/downloads");
    });

    describe("Requesting a report", () => {
        it("Successfully requests a report", () => {
            cy.get('[data-cy="report-type"]').select('Journal');
            cy.get('[data-cy="journal-types"]').select('ReceiptTransactions');
            cy.get('#date').type("2024-01-01");
            cy.get('#email').type("tina@test.com");
            cy.get('.govuk-button').contains('Download report').click();
            cy.url().should("include","/downloads?success=request_report&report_type=Receipt%20Transactions");
            cy.get('.moj-success-banner').contains('Your Receipt Transactions report is being prepared');
        });
    });
});
