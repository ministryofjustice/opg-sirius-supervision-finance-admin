describe("Downloading Files", () => {
    it("downloads a file by visiting the emailed link", () => {
        cy.visit("/download?uid=dGVzdC5jc3Y=");
        cy.contains(".govuk-heading-m", "test.csv is ready to download");
        cy.contains("button", "Download").click();
        cy.readFile("cypress/downloads/test.csv").should("exist");
    });
});
