describe("Downloading Files", () => {
    it("downloads a file by visiting the emailed link", () => {
        cy.visit("/download?uid=eyJLZXkiOiJ0ZXN0LmNzdiIsIlZlcnNpb25JZCI6InZwckF4c1l0TFZzYjVQOUhfcUhlTlVpVTlNQm5QTmN6In0=");
        cy.contains(".govuk-heading-m", "test.csv is ready to download");
        cy.contains("button", "Download").click();
        cy.readFile("cypress/downloads/test.csv", {timeout: 2000}).should("exist");
    });
});
