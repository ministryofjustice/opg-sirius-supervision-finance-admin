const { defineConfig } = require("cypress");
const AWS = require("aws-sdk");
const localstackEndpoint = "http://localstack:4566";

module.exports = defineConfig({
    fixturesFolder: false,
    e2e: {
        setupNodeEvents(on, config) {
            on("task", {
                log(message) {
                    console.log(message);

                    return null
                },
                table(message) {
                    console.table(message);

                    return null
                },
                uploadFileToS3({ bucket, key, body }) {
                    const s3 = new AWS.S3({
                        endpoint: localstackEndpoint,
                        s3ForcePathStyle: true,
                        accessKeyId: "test",
                        secretAccessKey: "test",
                    });

                    return s3.putObject({ Bucket: bucket, Key: key, Body: body }).promise()
                        .then(() => s3.headObject({ Bucket: bucket, Key: key }).promise())
                        .then(data => data.VersionId);
                },
                failed: require("cypress-failed-log/src/failed")()
            });
        },
        baseUrl: "http://localhost:8887/finance-admin",
        downloadsFolder: "cypress/downloads",
        specPattern: "cypress/e2e/**/*.cy.{js,ts}",
        screenshotsFolder: "cypress/screenshots",
        supportFile: "cypress/support/e2e.ts",
        modifyObstructiveCode: false,
    },
    viewportWidth: 1000,
    viewportHeight: 1000,
});