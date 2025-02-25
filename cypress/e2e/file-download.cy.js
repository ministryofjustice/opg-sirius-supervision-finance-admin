describe.skip("Downloading Files", () => {
    it("downloads a file by visiting the emailed link", () => {
        const bucket = 'opg-backoffice-finance-reports-local';
        const key = 'cypress-test.csv';
        const body = 'file content';

        cy.task('uploadFileToS3', { bucket, key, body }).then(versionId => {
            const uid = btoa(JSON.stringify({ Key: key, VersionId: versionId }));

            cy.visit(`/download?uid=${uid}`);
            cy.contains(".govuk-heading-m", `${key} is ready to download`);
            cy.contains("a", "Download").click();
            cy.readFile(`cypress/downloads/${key}`).should("exist");
        });
    });

    it("displays an error message when the file is not found", () => {
        const uid = btoa(JSON.stringify({ Key: 'non-existent-file.csv', VersionId: '123' }));

        cy.visit(`/download?uid=${uid}`);
        cy.contains(".govuk-heading-m", "Sorry, this link has expired. Please request a new report.");
        cy.contains("a", "Download").should("not.exist");
    });
});
