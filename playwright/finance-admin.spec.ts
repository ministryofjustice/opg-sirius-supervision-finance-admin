// playwright/e2e/finance_admin.spec.ts
import { test, expect, Page } from '@playwright/test';

test.describe('Finance Admin', () => {
    test.beforeEach(async ({ page }: { page: Page }) => {
        await page.goto('./uploads');
    });

    test.describe('Tabs', () => {
        [
            { title: 'Annual Invoicing Letters', path: 'annual-invoicing-letters' },
            { title: 'Uploads', path: 'uploads' },
            { title: 'Downloads', path: 'downloads' },
        ].forEach(({ title, path }) => {
            // You can also do it with test.describe() or with multiple tests as long the test name is unique.
            test(`navigates to ${title} correctly`, async ({ page }) => {
                const link = page.locator('.moj-sub-navigation__link', { hasText: title })
                await link.click();
                const regex = new RegExp(`.*${path}`);
                await expect(page).toHaveURL(regex);
                await expect(link)
                    .toHaveAttribute('aria-current', 'page');
            });
        });
    });
});