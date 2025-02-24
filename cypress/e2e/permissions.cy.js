const financeUser = "2"
const financeReporting = "3"
const corporateFinance = "4"

describe("Role-based permissions", () => {
    it("checks access for invalid roles", () => {
        cy.setCookie("x-test-user-id", financeUser);
        cy.visit("/uploads");
        //TODO: Add assertions
        cy.visit("/downloads");
        //TODO: Add assertions
    });

    it("checks permissions for Finance Reporting role", () => {
        cy.setCookie("x-test-user-id", financeReporting);
        cy.visit("/uploads");
        //TODO: Add assertions
        cy.visit("/downloads");
        //TODO: Add assertions
    });

    it("checks permissions for Corporate Finance role", () => {
        cy.setCookie("x-test-user-id", corporateFinance);
        cy.visit("/uploads");
        //TODO: Add assertions
        cy.visit("/downloads");
        //TODO: Add assertions
    });
});
