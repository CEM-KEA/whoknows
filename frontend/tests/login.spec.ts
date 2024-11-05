import test, { expect } from "@playwright/test";
import { config } from "dotenv";

config();

const baseUrl = process.env.TEST_FRONTEND_URL;
const testUserUsername = process.env.TEST_LOGIN_USERNAME ?? "";
const testUserPassword = process.env.TEST_LOGIN_PASSWORD ?? "";

test("can login", async ({ page }) => {
  await page.goto(`${baseUrl}/login`, { timeout: 10000 });
  await page.waitForSelector("#login-button", { state: "visible" });

  await page.fill("#login-username", testUserUsername);
  await page.fill("#login-password", testUserPassword);
  await page.click("#login-button");

  await page.waitForTimeout(5000);

  // if login succeeded it will either redirect to the home page or the change password page if the password needs to be changed
  // if login failed it will stay on the login page
  expect(page.url()).not.toBe(`${baseUrl}/login`);
});
