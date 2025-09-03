describe("Footer", () => {
    beforeEach(() => {
        cy.loginAs("Finance User Testing");
        cy.visit("/uploads");
    });

    it("should show the accessibility link", () => {
        cy.get('[data-cy="accessibilityStatement"]').should("contain", "Accessibility statement");
    });
});
