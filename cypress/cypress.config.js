const { defineConfig } = require("cypress");
const {
    S3Client,
    PutObjectCommand,
    HeadObjectCommand,
} = require("@aws-sdk/client-s3");
const localstackEndpoint = "http://localstack:4566";

module.exports = defineConfig({
    fixturesFolder: false,
    e2e: {
        setupNodeEvents(on) {
            on("task", {
                log(message) {
                    console.log(message);

                    return null;
                },
                table(message) {
                    console.table(message);

                    return null;
                },
                uploadFileToS3({ bucket, key, body }) {
                    const s3 = new S3Client({
                        region: "eu-west-1",
                        endpoint: localstackEndpoint,
                        forcePathStyle: true,
                        credentials: {
                            accessKeyId: "test",
                            secretAccessKey: "test",
                        },
                    });

                    return s3.send(new PutObjectCommand({ Bucket: bucket, Key: key, Body: body }))
                        .then(() => s3.send(new HeadObjectCommand({ Bucket: bucket, Key: key })))
                        .then((data) => data.VersionId);
                },
                failed: require("cypress-failed-log/src/failed")(),
            });
        },
        baseUrl: "http://localhost:8889/finance-admin",
        downloadsFolder: "cypress/downloads",
        specPattern: "cypress/e2e/**/*.cy.{js,ts}",
        screenshotsFolder: "cypress/screenshots",
        supportFile: "cypress/support/e2e.ts",
        modifyObstructiveCode: false,
    },
    viewportWidth: 1000,
    viewportHeight: 1000,
});
