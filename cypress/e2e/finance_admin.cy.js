describe("Finance Admin", () => {
    beforeEach(() => {
        cy.visit("/uploads");
    });

    describe("Tabs", () => {
        it("navigates between tabs correctly", () => {
            cy.get('[data-cy="annual-invoicing-letters-tab"]').click();
            cy.url().should("contain", "annual-invoicing-letters");
            // cy.contains(".moj-sub-navigation__link", "Annual Invoicing Letters")
            //     .should("have.attr", "aria-current", "page");

            cy.get('[data-cy="uploads-tab"]').click();
            cy.url().should("contain", "uploads");
            // cy.contains(".moj-sub-navigation__link", "Uploads")
            //     .should("have.attr", "aria-current", "page");


            cy.get('[data-cy="downloads-tab"]').click();
            cy.url().should("contain", "downloads");
            // cy.contains(".moj-sub-navigation__link", "Downloads")
            //     .should("have.attr", "aria-current", "page");
        });
    });
});
