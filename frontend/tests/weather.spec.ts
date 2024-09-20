import { test, expect } from "@playwright/test";
import { config } from "dotenv";
config();

const baseUrl = process.env.TEST_FRONTEND_URL;

// generate playwright tests to test that weather page gets weather content from backend
test("weather page gets weather content from backend", async ({ page }) => {
  // set a timeout of 10 seconds(max time to wait for the page to load)
  await page.goto(baseUrl + "/weather", { timeout: 10000 });

  // Check if the weather content is loaded
  const weatherContent = page.locator("#weather-content");
  await page.waitForTimeout(5000); // wait for 5 seconds
  await expect(weatherContent).toBeVisible();
});
