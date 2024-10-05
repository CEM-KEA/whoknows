import test, { expect } from "@playwright/test";
import { config } from "dotenv";

config();

const baseUrl = process.env.TEST_FRONTEND_URL;
const testUserUsername = process.env.TEST_LOGIN_USERNAME ?? "";
const testUserPassword = process.env.TEST_LOGIN_PASSWORD ?? "";

test("can login", async ({ page }) => {
  // Go to the login page and wait for the login button to appear
  await page.goto(`${baseUrl}/login`, { timeout: 10000 });
  await page.waitForSelector("#login-button", { state: "visible" });

  // Fill in the form
  await page.fill("#login-username", testUserUsername);
  await page.fill("#login-password", testUserPassword);

  // Click the login button and wait for the navigation
  await page.click("#login-button");

  // Add extra debug info after login
  await page.waitForTimeout(5000); // Give the page time to perform actions
  console.log("Current page URL after login click:", page.url());

  // Check if the navigation is working as expected
  await page.waitForURL(`${baseUrl}/`, { timeout: 10000 });
  
  console.log("Final page URL:", page.url());

  // Verify that we successfully navigated to the expected page
  expect(page.url()).toBe(`${baseUrl}/`);
});

test("can log in and then log out", async ({ page }) => {
  // set a timeout of 10 seconds(max time to wait for the page to load)
  await page.goto(`${baseUrl}/login`, { timeout: 10000 });

  // fill in the form and submit it.
  await page.fill("#login-username", testUserUsername);
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