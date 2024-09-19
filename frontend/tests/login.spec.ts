import test, { expect } from "@playwright/test";
import { config } from "dotenv";

config();

const baseUrl = process.env.TEST_FRONTEND_URL;
const testUserEmail = process.env.TEST_LOGIN_EMAIL ?? "";
const testUserPassword = process.env.TEST_LOGIN_PASSWORD ?? "";

test("can login", async ({ page }) => {
  // set a timeout of 10 seconds(max time to wait for the page to load)
  await page.goto(`${baseUrl}/login`, { timeout: 10000 });

  // fill in form
  await page.fill("#login-email", testUserEmail);
  await page.fill("#login-password", testUserPassword);
  await page.click("#login-button");

  // wait 5 second
  await page.waitForTimeout(5000);

  // expect to be redirected to the search page.
  expect(page.url()).toBe(`${baseUrl}/`);
  const logoutButton = await page.$("#login-logout-nav");
  expect(await logoutButton?.innerText()).toBe("Log out");
});

test("can log in and then log out", async ({ page }) => {
  // set a timeout of 10 seconds(max time to wait for the page to load)
  await page.goto(`${baseUrl}/login`, { timeout: 10000 });

  // fill in the form and submit it.
  await page.fill("#login-email", testUserEmail);
  await page.fill("#login-password", testUserPassword);
  await page.click("#login-button");

  // Wait for 5 seconds
  await page.waitForTimeout(5000);

  // Expect to be redirected to the search page.
  expect(page.url()).toBe(`${baseUrl}/`);

  // Click the logout button
  await page.click("#login-logout-nav");

  // Wait for 5 seconds
  await page.waitForTimeout(5000);

  // Expect to be redirected to the login page.
  expect(page.url()).toBe(`${baseUrl}/login`);
  const loginButton = await page.$("#login-logout-nav");
  expect(await loginButton?.innerText()).toBe("Log in");
});
