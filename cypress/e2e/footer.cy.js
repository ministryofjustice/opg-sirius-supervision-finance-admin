describe("Footer", () => {
    beforeEach(() => {
        cy.visit("/uploads");
    });

    it("should show the accessibility link", () => {
        cy.get('[data-cy="accessibilityStatement"]').should("contain", "Accessibility statement");
    });
});
