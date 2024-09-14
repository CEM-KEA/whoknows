import test, { expect } from "@playwright/test";
import {config} from "dotenv";

config();

const baseUrl = process.env.TEST_FRONTEND_URL;

// Test search functionality, playwright
test('search', async ({ page }) => {
  await page.goto(`${baseUrl}`);

  // Search for 'Test'.
  await page.fill('#search', 'Test');
  await page.keyboard.press('Enter');

  // Wait for 1 seconds
  await page.waitForTimeout(1000);
    
  // Expect content to be shown.
  const searchResults = page.locator('#search-results>li');
  const count = await searchResults.count();
  expect(count).toBeGreaterThan(0);
});
