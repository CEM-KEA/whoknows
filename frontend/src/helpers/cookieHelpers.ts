import Cookies from "universal-cookie";

const cookies = new Cookies();

export function getJWTTokenFromCookies(): string {
  return cookies.get("jwt_authorization");
}

export function setJWTTokenInCookies(token: string): void {
  cookies.set("jwt_authorization", token);
}

export function removeJWTTokenFromCookies(): void {
  cookies.remove("jwt_authorization");
}
