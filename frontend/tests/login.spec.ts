import test, { expect } from "@playwright/test";
import { config } from "dotenv";

config();

const baseUrl = process.env.TEST_FRONTEND_URL;
//test login playwright
test("login", async ({ page }) => {
  await page.goto(`${baseUrl}/login`);

  // Fill in the form and submit it.
  await page.fill("#login-username", "testuser");
  await page.fill("#login-password", "password");
  await page.click("#login-button");

  // Wait for 1 seconds
  await page.waitForTimeout(1000);

  // Expect to be redirected to the search page.
  expect(page.url()).toBe(`${baseUrl}/`);
});
