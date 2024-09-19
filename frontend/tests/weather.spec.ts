import { test, expect } from "@playwright/test";
import { config } from "dotenv";
config();

const baseUrl = process.env.TEST_FRONTEND_URL;

// generate playwright tests to test that weather page gets weather content from backend
test("weather page gets weather content from backend", async ({ page }) => {
  await page.goto(baseUrl + "/weather");

  // Check if the weather content is loaded
  const weatherContent = page.locator("#weather-content");
  await page.waitForTimeout(1000); // wait for 1 second
  await expect(weatherContent).toBeVisible();
});
