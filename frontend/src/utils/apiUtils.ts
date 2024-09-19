import Cookies from "universal-cookie";

const apiUrl = import.meta.env.VITE_API_URL;
const cookies = new Cookies();

/**
 * Sends a GET request to the API, url is the path to the endpoint and should start with a /.
 *
 * Example: apiGet("/users") will send a GET request to /api/users
 */
export async function apiGet<TResBody>(url: string, requireAuth?: boolean): Promise<TResBody> {
  const res = await fetch(apiUrl + url, {
    headers: {
      Authorization: requireAuth ? `Bearer ${cookies.get("jwt_authorization")}` : ""
    }
  });
  if (!res.ok) {
    throw new Error(res.statusText);
  }
  return await res.json();
}

/**
 * Sends a POST request to the API, url is the path to the endpoint and should start with a /.
 *
 * Example: apiPost("/users", {name: "John Doe"}) will send a POST request to /api/users with the body {name: "John Doe"}
 */
export async function apiPost<TReqBody, TResBody>(
  url: string,
  data: TReqBody,
  requireAuth?: boolean
): Promise<TResBody> {
  const res = await fetch(apiUrl + url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: requireAuth ? `Bearer ${cookies.get("jwt_authorization")}` : ""
    },
    body: JSON.stringify(data)
  });
  if (!res.ok) {
    throw new Error(res.statusText);
  }
  return await res.json();
}

/**
 * Sends a PUT request to the API, url is the path to the endpoint and should start with a /.
 *
 * Example: apiPut("/users/1", {name: "Jane Doe"}) will send a PUT request to /api/users/1 with the body {name: "Jane Doe"}
 */
export async function apiPut<TReqBody, TResBody>(
  url: string,
  data: TReqBody,
  requireAuth?: boolean
): Promise<TResBody> {
  const res = await fetch(apiUrl + url, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
      Authorization: requireAuth ? `Bearer ${cookies.get("jwt_authorization")}` : ""
    },
    body: JSON.stringify(data)
  });
  if (!res.ok) {
    throw new Error(res.statusText);
  }
  return await res.json();
}

/**
 * Sends a DELETE request to the API, url is the path to the endpoint and should start with a /.
 *
 * Example: apiDelete("/users/1") will send a DELETE request to /api/users/1
 */
export async function apiDelete<TResBody>(url: string, requireAuth?: boolean): Promise<TResBody> {
  const res = await fetch(apiUrl + url, {
    method: "DELETE",
    headers: {
      Authorization: requireAuth ? `Bearer ${cookies.get("jwt_authorization")}` : ""
    }
  });
  if (!res.ok) {
    throw new Error(res.statusText);
  }
  return await res.json();
}
