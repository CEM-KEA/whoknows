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
  await page.waitForURL(`${baseUrl}/`, { timeout: 10000 });

  expect(page.url()).toBe(`${baseUrl}/`);
});

test("can log in and then log out", async ({ page }) => {
  await page.goto(`${baseUrl}/login`, { timeout: 10000 });

  await page.fill("#login-username", testUserUsername);
  await page.fill("#login-password", testUserPassword);
  await page.click("#login-button");

  await page.waitForTimeout(5000);

  expect(page.url()).toBe(`${baseUrl}/`);

  await page.click("#login-logout-nav");

  await page.waitForTimeout(5000);

  // Expect to be redirected to the login page.
  expect(page.url()).toBe(`${baseUrl}/login`);
  const loginButton = await page.$("#login-logout-nav");
  expect(await loginButton?.innerText()).toBe("Log in");
});