describe("Uploading Files", () => {
    beforeEach(() => {
        cy.visit("/uploads");
    });

    describe("Upload file", () => {
        it("Uploads file successfully", () => {
            cy.get('[data-cy=\"upload-type\"]').select('PAYMENTS_MOTO_CARD');
            cy.get('#file-upload').selectFile('cypress/fixtures/feemoto_01102024normal.csv');
            cy.get('#upload-date').type("2024-10-01");
            cy.get('.govuk-button').contains('Upload file').click();
            cy.url().should("include","/uploads?success=upload");
            cy.get('.moj-banner').contains('File successfully uploaded');
        });

        it("Validates missing upload date", () => {
            cy.get('[data-cy=\"upload-type\"]').select('PAYMENTS_MOTO_CARD');
            cy.get('#file-upload').selectFile('cypress/fixtures/feemoto_01102024normal.csv');
            cy.get('.govuk-button').contains('Upload file').click();
            cy.get('#f-UploadDate').contains('Please enter a date');
            cy.get('#f-UploadDate').should('have.class', 'govuk-form-group--error');
        });

        it("Validates missing file", () => {
            cy.get('[data-cy=\"upload-type\"]').select('PAYMENTS_MOTO_CARD');
            cy.get('.govuk-button').contains('Upload file').click();
            cy.get('.govuk-error-summary').contains('No file uploaded');
            cy.get('#f-FileUpload').contains('No file uploaded');
            cy.get('#f-FileUpload').should('have.class', 'govuk-form-group--error');
        });

        it("Validates empty headers", () => {
            cy.get('[data-cy=\"upload-type\"]').select('PAYMENTS_MOTO_CARD');
            cy.get('#file-upload').selectFile('cypress/fixtures/feemoto_01122024normal.csv');
            cy.get('#upload-date').type("2024-12-01");
            cy.get('.govuk-button').contains('Upload file').click();
            cy.get('.govuk-error-summary').contains('Failed to read CSV headers');
            cy.get('#f-FileUpload').contains('Failed to read CSV headers');
            cy.get('#f-FileUpload').should('have.class', 'govuk-form-group--error');
        });

        it("Validates CSV headers", () => {
            cy.get('[data-cy=\"upload-type\"]').select('PAYMENTS_MOTO_CARD');
            cy.get('#file-upload').selectFile('cypress/fixtures/feemoto_02102024normal.csv');
            cy.get('#upload-date').type("2024-10-02");
            cy.get('.govuk-button').contains('Upload file').click();
            cy.get('.govuk-error-summary').contains('CSV headers do not match for the file being uploaded');
            cy.get('#f-FileUpload').contains('CSV headers do not match for the file being uploaded');
            cy.get('#f-FileUpload').should('have.class', 'govuk-form-group--error');
        });
    });

    describe("Form fields", () => {
        const selections = {
            "PAYMENTS_MOTO_CARD": ['#file-upload', '#upload-date', '#email-field'],
            "PAYMENTS_ONLINE_CARD": ['#file-upload', '#upload-date', '#email-field'],
            "PAYMENTS_OPG_BACS": ['#file-upload', '#upload-date', '#email-field'],
            "PAYMENTS_SUPERVISION_BACS": ['#file-upload', '#upload-date', '#email-field'],
            "PAYMENTS_SUPERVISION_CHEQUE": ['#file-upload', '#upload-date', '#pis-number', '#email-field'],
            "DEBT_CHASE": ['#file-upload', '#email-field'],
            "DEPUTY_SCHEDULE": ['#file-upload', '#email-field'],
            "DIRECT_DEBITS_COLLECTIONS": ['#file-upload', '#upload-date', '#email-field'],
            "MISAPPLIED_PAYMENTS": ['#file-upload', '#email-field'],
        }

        it("Displays the correct fields for selection", () => {
            for (const [upload, fields] of Object.entries(selections)) {
                cy.get('[data-cy=\"upload-type\"]').select(upload);
                fields.forEach((field) => {
                   cy.get(field).should('be.visible');
                });
            }
        });
    });
});
