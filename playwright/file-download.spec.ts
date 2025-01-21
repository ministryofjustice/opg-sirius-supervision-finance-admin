// playwright/e2e/file-download.spec.js
const { test, expect } = require('@playwright/test');
import { S3Client, PutObjectCommand } from "@aws-sdk/client-s3";

const s3 = new S3Client({
    endpoint: process.env.S3_ENDPOINT,
    region: "eu-west-1",
    forcePathStyle: true,
    credentials: {
        accessKeyId: "test",
        secretAccessKey: "test",
    },
});

test.describe('Downloading Files', () => {
    test('downloads a file by visiting the emailed link', async ({ page }) => {
        const bucket = 'opg-backoffice-finance-reports-local';
        const key = 'cypress-test.csv';
        const body = 'file content';

        const command = new PutObjectCommand({
            Bucket: bucket,
            Key: key,
            Body: body,
        });

        const uploadResult = await s3.send(command)

        const versionId = uploadResult.VersionId;
        const uid = Buffer.from(JSON.stringify({ Key: key, VersionId: versionId })).toString('base64');

        // Visit the download link
        await page.goto(`./download?uid=${uid}`);
        await expect(page.locator('.govuk-heading-m')).toHaveText(`${key} is ready to download`);
        await page.getByRole('button', { name: 'Upload' }).click();

        // Verify the file is downloaded
        const downloadPath = await page.context().downloadsPath();
        const filePath = `${downloadPath}/${key}`;
        const fs = require('fs');
        expect(fs.existsSync(filePath)).toBeTruthy();
    });

    test('displays an error message when the file is not found', async ({ page }) => {
        const uid = Buffer.from(JSON.stringify({ Key: 'non-existent-file.csv', VersionId: '123' })).toString('base64');

        // Visit the download link
        await page.goto(`./download?uid=${uid}`);
        await expect(page.locator('.govuk-heading-m')).toHaveText('Sorry, this link has expired. Please request a new report.');
        await expect(page.getByRole('button', { name: 'Download' })).not.toBeVisible();
    });
});