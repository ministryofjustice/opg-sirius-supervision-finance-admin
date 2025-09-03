describe("Downloads", () => {
    beforeEach(() => {
        cy.loginAs("Finance User Testing");
        cy.visit("/downloads");
    });

    describe("Requesting a report", () => {
        it("Successfully requests a report", () => {
            cy.get('[data-cy="report-type"]').select('Journal');
            cy.get('[data-cy="journal-types"]').select('ReceiptTransactions');
            cy.get('#date').type("2024-01-01");
            cy.get('#email').type("tina@test.com");
            cy.get('.govuk-button').contains('Download report').click();
            cy.url().should("include", "/downloads?success=request_report&report_type=Receipt%20Transactions");
            cy.get('.moj-success-banner').contains('Your Receipt Transactions report is being prepared');
        });
    });


    describe("Form fields", () => {
        const selections = {
            "Journal": ["#journal-types", "#date", "#email"],
            "Schedule": ["#schedule-types", "#date", "#email"],
            "Debt": ["#debt-types", "#email"],
        }

        it("Displays the correct fields for non-account-receivable selection", () => {
            for (const [upload, fields] of Object.entries(selections)) {
                cy.get('[data-cy=\"report-type\"]').select(upload);
                fields.forEach((field) => {
                    cy.get(field).should('be.visible');
                });
            }
        });

        const ARSelections = {
            "AgedDebt": ["#date-from", "#date-to", "#email"],
            "AgedDebtByCustomer": ["#email"],
            "UnappliedReceipts": ["#date-from", "#date-to", "#email"],
            "ARPaidInvoice": ["#date-from", "#date-to", "#email"],
            "TotalReceipts": ["#date-from", "#date-to", "#email"],
            "BadDebtWriteOff": ["#date-from", "#date-to", "#email"],
            "FeeAccrual": ["#email"],
            "InvoiceAdjustments": ["#date-from", "#date-to", "#email"],
        }

        it("Displays the correct fields for account-receivable selection", () => {
            cy.get('[data-cy=\"report-type\"]').select("AccountsReceivable");
            for (const [report, fields] of Object.entries(ARSelections)) {
                cy.get('#account-types').select(report);
                fields.forEach((field) => {
                    cy.get(field).should('be.visible');
                });
            }
        });
    });
});
