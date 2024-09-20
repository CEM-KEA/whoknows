import test, { expect } from "@playwright/test";
import { config } from "dotenv";

config();

const baseUrl = process.env.TEST_FRONTEND_URL;

// Test search functionality, playwright
test("search", async ({ page }) => {
  // Set a timeout of 10 seconds(max time to wait for the page to load)
  await page.goto(`${baseUrl}`, { timeout: 10000 });

  // Search for 'Test'.
  await page.fill("#search", "Test");
  await page.keyboard.press("Enter");

  // Wait for 5 seconds
  await page.waitForTimeout(5000);

  // Expect content to be shown.
  const searchResults = page.locator("#search-results>li");
  const count = await searchResults.count();
  expect(count).toBeGreaterThan(0);
});
