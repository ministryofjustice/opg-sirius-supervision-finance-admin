describe("Finance Admin", () => {
    beforeEach(() => {
        cy.visit("/uploads");
    });

    describe("Upload file", () => {
        it("Uploads file successfully", () => {
            cy.get('[data-cy=\"report-upload-type\"]').select('PAYMENTS_MOTO_CARD')
            cy.get('#file-upload').selectFile('cypress/fixtures/feemoto_01:10:2024normal.csv')
            cy.get('#upload-date').type("2024-10-01");
            cy.get('.govuk-button').contains('Upload report').click()
            cy.url().should("include","/uploads?success=upload");
            cy.get('.moj-banner').contains('File successfully uploaded')
        });

        it("Validates missing file", () => {
            cy.get('[data-cy=\"report-upload-type\"]').select('PAYMENTS_MOTO_CARD')
            cy.get('.govuk-button').contains('Upload report').click()
            cy.get('.govuk-error-summary').contains('No file uploaded')
            cy.get('#f-FileUpload').contains('No file uploaded')
            cy.get('#f-FileUpload').should('have.class', 'govuk-form-group--error')
        })

        it("Validates empty headers", () => {
            cy.get('[data-cy=\"report-upload-type\"]').select('DEBT_CHASE')
            cy.get('#file-upload').selectFile('cypress/fixtures/empty_report.csv')
            cy.get('.govuk-button').contains('Upload report').click()
            cy.get('.govuk-error-summary').contains('Failed to read CSV headers')
            cy.get('#f-FileUpload').contains('Failed to read CSV headers')
            cy.get('#f-FileUpload').should('have.class', 'govuk-form-group--error')
        })

        it("Validates CSV headers", () => {
            cy.get('[data-cy=\"report-upload-type\"]').select('PAYMENTS_MOTO_CARD');
            cy.get('#file-upload').selectFile('cypress/fixtures/feemoto_02:10:2024normal.csv')
            cy.get('#upload-date').type("2024-10-02");
            cy.get('.govuk-button').contains('Upload report').click()
            cy.get('.govuk-error-summary').contains('CSV headers do not match for the report trying to be uploaded')
            cy.get('#f-FileUpload').contains('CSV headers do not match for the report trying to be uploaded')
            cy.get('#f-FileUpload').should('have.class', 'govuk-form-group--error')
        });
    });
});
