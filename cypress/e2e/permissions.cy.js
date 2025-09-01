const allAccess = "1"
const financeReportingOnly = "2"
const corporateFinanceOnly = "3"

describe("Role-based permissions", () => {
    it("checks permissions for all access", () => {
        cy.setCookie("x-test-user-id", allAccess);
        cy.visit("/uploads");
        cy.get('[data-cy="upload-type"]').find("option")
            .should("have.length", 14);

        cy.visit("/downloads");
        cy.contains(".govuk-heading-m", "Download a report");
    });

    it("checks permissions for Finance Reporting only", () => {
        cy.setCookie("x-test-user-id", financeReportingOnly);
        cy.visit("/uploads");
        cy.get('[data-cy="upload-type"]').find("option")
            .should("have.length", 4);

        cy.visit("/downloads");
        cy.contains(".govuk-heading-m", "Download a report");
    });

    it("checks permissions for Corporate Finance only", () => {
        // should see nothing as Finance Reporting required to access page
        cy.setCookie("x-test-user-id", corporateFinanceOnly);
        cy.visit("/uploads", {failOnStatusCode: false});
        cy.contains(".govuk-heading-l", "Forbidden");

        cy.visit("/downloads", {failOnStatusCode: false});
        cy.contains(".govuk-heading-l", "Forbidden");
    });
});
