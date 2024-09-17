import Cookies from "universal-cookie";

const cookies = new Cookies();

export function getJWTTokenFromCookies(): string {
  if (!getUserCookieAccept()) return "";
  return cookies.get("jwt_authorization");
}

export function setJWTTokenInCookies(token: string): void {
  if (!getUserCookieAccept()) return;
  cookies.set("jwt_authorization", token);
}

export function removeJWTTokenFromCookies(): void {
  if (!getUserCookieAccept()) return;
  cookies.remove("jwt_authorization");
}

export function setUserCookieAccept(choice: boolean): void {
  localStorage.setItem("user_cookie_accept", choice.toString());
}

export function getUserCookieAccept(): boolean | null {
  const accept = localStorage.getItem("user_cookie_accept");
  if (!accept) {
    return null;
  }
  return accept === "true";
}
