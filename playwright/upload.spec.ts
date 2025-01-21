// playwright/e2e/upload.spec.ts
import { test, expect, Page } from '@playwright/test';
// @ts-ignore
import path from 'path';

test.describe('Uploading Files', () => {
    test.beforeEach(async ({ page }: { page: Page }) => {
        await page.goto('./uploads');
    });

    test.describe('Upload file', () => {
        test('Uploads file successfully', async ({ page }: { page: Page }) => {
            await page.getByLabel('Select a report type').selectOption('PAYMENTS_MOTO_CARD');
            const filePath = path.resolve(__dirname, './fixtures/feemoto_01:10:2024normal.csv');
            await page.getByLabel('Date').fill('2024-10-01');
            await page.getByLabel('Upload a file').setInputFiles(filePath);
            await page.getByRole('button', { name: 'Upload report' }).click();
            await expect(page).toHaveURL(/\/uploads\?success=upload/);
            await expect(page.locator('.moj-banner')).toHaveText(/File successfully uploaded/);
        });

        test('Validates missing file', async ({ page }: { page: Page }) => {
            await page.getByLabel('Select a report type').selectOption('PAYMENTS_MOTO_CARD');
            await page.getByRole('button', { name: 'Upload report' }).click();
            await expect(page.locator('.govuk-error-summary')).toHaveText(/No file uploaded/);
            await expect(page.locator('#f-FileUpload')).toHaveText(/No file uploaded/);
            await expect(page.locator('#f-FileUpload')).toHaveClass(/govuk-form-group--error/);
        });

        test('Validates empty headers', async ({ page }: { page: Page }) => {
            await page.getByLabel('Select a report type').selectOption('DEBT_CHASE');
            const filePath = path.resolve(__dirname, './fixtures/empty_report.csv');
            await page.getByLabel('Upload a file').setInputFiles(filePath);
            await page.getByRole('button', { name: 'Upload report' }).click();
            await expect(page.locator('.govuk-error-summary')).toHaveText(/Failed to read CSV headers/);
            await expect(page.locator('#f-FileUpload')).toHaveText(/Failed to read CSV headers/);
            await expect(page.locator('#f-FileUpload')).toHaveClass(/govuk-form-group--error/);
        });

        test('Validates CSV headers', async ({ page }: { page: Page }) => {
            await page.getByLabel('Select a report type').selectOption('PAYMENTS_MOTO_CARD');
            const filePath = path.resolve(__dirname, './fixtures/feemoto_02:10:2024normal.csv');
            await page.getByLabel('Date').fill('2024-10-02');
            await page.getByLabel('Upload a file').setInputFiles(filePath);
            await page.getByRole('button', { name: 'Upload report' }).click();
            await expect(page.locator('.govuk-error-summary')).toHaveText(/CSV headers do not match for the report trying to be uploaded/);
            await expect(page.locator('#f-FileUpload')).toHaveText(/CSV headers do not match for the report trying to be uploaded/);
            await expect(page.locator('#f-FileUpload')).toHaveClass(/govuk-form-group--error/);
        });
    });
});