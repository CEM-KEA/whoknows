export interface IAuthResponse {
  statusCode: number;
  message: string;
}

export interface ILoginRequest {
  username: string;
  password: string;
}

export interface ILoginSession {
  username: string;
}

export interface ILoginResponse {
  token: string;
}
