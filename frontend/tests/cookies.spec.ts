import test, { expect } from "@playwright/test";
import * as cookies from "../src/helpers/cookieHelpers";

test("set, get and remove cookies", async () => {
  const token = "test_token";
  cookies.setJWTTokenInCookies(token);
  expect(cookies.getJWTTokenFromCookies()).toBe(token);
  cookies.removeJWTTokenFromCookies();
  expect(cookies.getJWTTokenFromCookies()).toBeUndefined();
});
