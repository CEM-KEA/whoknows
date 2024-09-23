import test, { expect } from "@playwright/test";
import { config } from "dotenv";

config();

const baseUrl = process.env.TEST_FRONTEND_URL;

test("Register button enabled for valid inputs", async ({ page }) => {
  await page.goto(`${baseUrl}/register`);

  await page.fill('input[name="username"]', "test");
  await page.fill('input[name="email"]', "test@test.com");
  await page.fill('input[name="password"]', "password");
  await page.fill('input[name="repeat-password"]', "password");

  // test that the submit button is not disabled
  const registerButton = page.locator("#register-button");
  expect(await registerButton?.isVisible()).toBe(true);
  expect(await registerButton?.isEnabled()).toBe(true);
});

test("Register button disabled for invalid username", async ({ page }) => {
  await page.goto(`${baseUrl}/register`);

  // username must be between 3 and 100 characters
  await page.fill('input[name="username"]', "ab");
  await page.fill('input[name="email"]', "test@test.com");
  await page.fill('input[name="password"]', "password");
  await page.fill('input[name="repeat-password"]', "password");

  // test that the submit button is disabled
  const registerButton = page.locator("#register-button");
  expect(await registerButton?.isVisible()).toBe(true);
  expect(await registerButton?.isDisabled()).toBe(true);
});

test("Register button disabled for invalid email", async ({ page }) => {
  await page.goto(`${baseUrl}/register`);

  await page.fill('input[name="username"]', "test");
  await page.fill('input[name="email"]', "test");
  await page.fill('input[name="password"]', "password");
  await page.fill('input[name="repeat-password"]', "password");

  // test that the submit button is disabled
  const registerButton = page.locator("#register-button");
  expect(await registerButton?.isVisible()).toBe(true);
  expect(await registerButton?.isDisabled()).toBe(true);
});

test("Register button disabled for invalid password", async ({ page }) => {
  await page.goto(`${baseUrl}/register`);

  await page.fill('input[name="username"]', "test");
  await page.fill('input[name="email"]', "test@test.com");
  // password must be at least 6 characters
  await page.fill('input[name="password"]', "pass");
  await page.fill('input[name="repeat-password"]', "pass");

  // test that the submit button is disabled
  const registerButton = page.locator("#register-button");
  expect(await registerButton?.isVisible()).toBe(true);
  expect(await registerButton?.isDisabled()).toBe(true);
});

test("Register button disabled for passwords that do not match", async ({ page }) => {
  await page.goto(`${baseUrl}/register`);

  await page.fill('input[name="username"]', "test");
  await page.fill('input[name="email"]', "test@test.com");
  await page.fill('input[name="password"]', "password");
  await page.fill('input[name="repeat-password"]', "password2");

  // test that the submit button is disabled
  const registerButton = page.locator("#register-button");
  expect(await registerButton?.isVisible()).toBe(true);
  expect(await registerButton?.isDisabled()).toBe(true);
});
