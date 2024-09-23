import test, { expect, Page } from "@playwright/test";
import { config } from "dotenv";

config();

const baseUrl = process.env.TEST_FRONTEND_URL;

const validUsername = "test";
const validEmail = "test@test.com";
const validPassword = "password";

const invalidUsername = "ab";
const invalidEmail = "test";
const invalidPassword = "pass";

async function fillRegisterForm(
  page: Page,
  username: string,
  email: string,
  password: string,
  repeatPassword: string
) {
  await page.fill('input[name="username"]', username);
  await page.fill('input[name="email"]', email);
  await page.fill('input[name="password"]', password);
  await page.fill('input[name="repeat-password"]', repeatPassword);
}

async function checkRegisterButtonState(page: Page, shouldBeEnabled: boolean) {
  const registerButton = page.locator("#register-button");
  expect(await registerButton?.isVisible()).toBe(true);
  if (shouldBeEnabled) {
    expect(await registerButton?.isEnabled()).toBe(true);
  } else {
    expect(await registerButton?.isDisabled()).toBe(true);
  }
}

test("Register button enabled for valid inputs", async ({ page }) => {
  await page.goto(`${baseUrl}/register`);
  await fillRegisterForm(page, validUsername, validEmail, validPassword, validPassword);
  await checkRegisterButtonState(page, true);
});

test("Register button disabled for invalid username", async ({ page }) => {
  await page.goto(`${baseUrl}/register`);
  await fillRegisterForm(page, invalidUsername, validEmail, validPassword, validPassword);
  await checkRegisterButtonState(page, false);
});

test("Register button disabled for invalid email", async ({ page }) => {
  await page.goto(`${baseUrl}/register`);
  await fillRegisterForm(page, validUsername, invalidEmail, validPassword, validPassword);
  await checkRegisterButtonState(page, false);
});

test("Register button disabled for invalid password", async ({ page }) => {
  await page.goto(`${baseUrl}/register`);
  await fillRegisterForm(page, validUsername, validEmail, invalidPassword, invalidPassword);
  await checkRegisterButtonState(page, false);
});

test("Register button disabled for mismatched passwords", async ({ page }) => {
  await page.goto(`${baseUrl}/register`);
  await fillRegisterForm(page, validUsername, validEmail, validPassword, invalidPassword);
  await checkRegisterButtonState(page, false);
});
